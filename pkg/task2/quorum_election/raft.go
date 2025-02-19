package quorum_election

import (
	"errors"
	"fmt"
	"graid-tech-assignment/pkg/task2/message"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	Follower = iota
	Candidate
	Leader
	Dead
)

func StateToString(state int) string {
	switch state {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	case Dead:
		return "Dead"
	default:
		return "Unknown"
	}
}

type RaftCluster struct {
	nodes *map[int]*Raft
}

type Raft struct {
	mu sync.Mutex // Lock to protect shared access to this peer's state
	ID int
	//peerIDs []int // peerIDs list the IDs of peers in the cluster
	peers map[int]*Raft

	currentTerm int
	votedFor    int

	state    int // 0 = follower, 1 = candidate, 2 = leader, 3 = dead
	leaderId int

	// Timer or Ticker
	electionResetEvent time.Time
}

func NewRaft(id int, ready <-chan any) *Raft {
	rf := &Raft{
		ID:          id,
		currentTerm: 0,
		votedFor:    -1,
		state:       0,
		leaderId:    -1,
		peers:       make(map[int]*Raft),
	}

	go func() {
		<-ready
		rf.mu.Lock()
		rf.electionResetEvent = time.Now()
		rf.mu.Unlock()
		rf.runElectionTimer()
	}()
	return rf
}

type RequestVoteInput struct {
	Term        int
	CandidateId int
}
type RequestVoteOutput struct {
	Term        int
	VoteGranted bool
}

type HeartbeatInput struct {
	Term     int
	LeaderId int
}
type HeartbeatOutput struct {
	Term    int
	Success bool
}

func (rf *Raft) Report() (id int, term int, state string) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.ID, rf.currentTerm, StateToString(rf.state)
}

func (rf *Raft) Stop() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.state = Dead
	for _, peer := range rf.peers {
		peer.removePeer(rf.ID)
	}
	log.Printf("server %v is being stopped", rf.ID)
}

func (rf *Raft) Start() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.state == Dead {
		rf.state = Follower
	}
}

func (rf *Raft) removePeer(peerId int) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	//newPeerIDs := make([]int, 0)
	//for _, val := range rf.peerIDs {
	//	if val != peerId {
	//		newPeerIDs = append(newPeerIDs, val)
	//	}
	//}
	//rf.peerIDs = newPeerIDs
	//delete(rf.peers, peerId)
	//log.Printf("[%d]peerIDs has been changes: %+v, %+v", rf.ID, rf.peerIDs, rf.peers)
}

func (rf *Raft) AppendPeer(peerId int, peer *Raft) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	for _, p := range rf.peers {
		if p.ID == peerId {
			return
		}
	}
	//rf.peerIDs = append(rf.peerIDs, peerId)
	rf.peers[peerId] = peer
	//log.Printf("[%d]peerIDs has been changes: %+v, %+v", rf.ID, rf.peerIDs, rf.peers)
}

func (rf *Raft) electionTimeout() time.Duration {
	return time.Duration(150+rand.Intn(150)) * time.Millisecond
}

// runElectionTimer is blocking, and should be run on a separate goroutine.
func (rf *Raft) runElectionTimer() {
	timeout := rf.electionTimeout()
	rf.mu.Lock()
	termStarted := rf.currentTerm
	rf.mu.Unlock()

	log.Printf("[%d]election timer started for term %d, timeout %v", rf.ID, termStarted, timeout)

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		<-ticker.C
		rf.mu.Lock()
		if rf.state != Candidate && rf.state != Follower {
			log.Printf("[%d]election timer state=%s in term %d, skipping", rf.ID, StateToString(rf.state), termStarted)
			rf.mu.Unlock()
			return
		}

		if termStarted != rf.currentTerm {
			log.Printf("[%d]election timer term changed from %d to %d, skipping", rf.ID, termStarted, rf.currentTerm)
			rf.mu.Unlock()
			return
		}

		if elapsed := time.Since(rf.electionResetEvent); elapsed >= timeout {
			log.Printf("[%d]failed heartbeat with Leader %d", rf.ID, rf.leaderId)
			rf.startElection()
			rf.mu.Unlock()
			return
		}
		rf.mu.Unlock()
	}

}

func (rf *Raft) sendRequestVote(input RequestVoteInput, peer int) (*RequestVoteOutput, error) {
	request := &message.RequestVoteMessage{
		SenderID:    rf.ID,
		PeerID:      peer,
		Term:        input.Term,
		CandidateId: input.CandidateId,
	}
	reply, err := rf.peers[peer].handleRequestVote(request)
	if err != nil {
		return nil, err
	}
	return &RequestVoteOutput{
		Term:        reply.Term,
		VoteGranted: reply.VoteGranted,
	}, nil
}

func (rf *Raft) handleRequestVote(input *message.RequestVoteMessage) (*message.RequestVoteMessageReply, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.state == Dead {
		return nil, errors.New(fmt.Sprintf("[%d]is dead", rf.ID))
	}

	if input.Term > rf.currentTerm {
		rf.becomeFollower(input.Term)
	}

	reply := &message.RequestVoteMessageReply{
		SenderID: rf.ID,
		PeerID:   input.GetSenderId(),
	}
	if input.Term == rf.currentTerm &&
		(rf.votedFor == -1 || rf.votedFor == input.CandidateId) {
		log.Printf("[%d]Accept member %d to be leader", rf.ID, input.CandidateId)
		reply.VoteGranted = true
		rf.votedFor = input.CandidateId
		rf.electionResetEvent = time.Now()
	} else {
		reply.VoteGranted = false
	}

	reply.Term = rf.currentTerm

	return reply, nil
}

func (rf *Raft) sendHeartbeats(input HeartbeatInput, peer int) (*HeartbeatOutput, error) {
	request := &message.HeartbeatMessage{
		SenderID: rf.ID,
		PeerID:   peer,
		Term:     input.Term,
		LeaderId: input.LeaderId,
	}
	reply, err := rf.peers[peer].handleHeartbeats(request)
	if err != nil {
		return nil, err
	}
	return &HeartbeatOutput{
		Term:    reply.Term,
		Success: reply.Success,
	}, nil
}

func (rf *Raft) handleHeartbeats(input *message.HeartbeatMessage) (*message.HeartbeatMessageReply, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.state == Dead {
		return nil, errors.New(fmt.Sprintf("[%d]is dead", rf.ID))
	}

	if input.Term > rf.currentTerm {
		rf.becomeFollower(input.Term)
	}

	reply := &message.HeartbeatMessageReply{
		SenderID: rf.ID,
		PeerID:   input.GetSenderId(),
		Success:  false,
	}
	if input.Term == rf.currentTerm {
		if rf.state != Follower {
			rf.becomeFollower(input.Term)
		}
		rf.electionResetEvent = time.Now()
		rf.leaderId = input.LeaderId
		reply.Success = true
	}

	reply.Term = rf.currentTerm

	return reply, nil
}

func (rf *Raft) becomeFollower(term int) {
	rf.state = Follower
	rf.currentTerm = term
	rf.votedFor = -1
	rf.electionResetEvent = time.Now()
	go rf.runElectionTimer()
}

func (rf *Raft) becomeLeader() {
	log.Printf("[%d]become leader", rf.ID)
	rf.state = Leader
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			rf.mu.Lock()
			if rf.state != Leader {
				rf.mu.Unlock()
				return
			}

			savedCurrentTerm := rf.currentTerm
			rf.mu.Unlock()

			for _, peer := range rf.peers {
				input := HeartbeatInput{
					Term:     savedCurrentTerm,
					LeaderId: rf.ID,
				}
				go func() {
					output, err := rf.sendHeartbeats(input, peer.ID)
					if err != nil {
						return
					}
					rf.mu.Lock()
					defer rf.mu.Unlock()

					if output.Term > savedCurrentTerm {
						rf.becomeFollower(output.Term)
					}
				}()
			}

		}
	}()
}

func (rf *Raft) startElection() {
	rf.state = Candidate
	rf.currentTerm++
	savedCurrentTerm := rf.currentTerm
	rf.electionResetEvent = time.Now()
	rf.votedFor = rf.ID
	log.Printf("[%d]I want to be leader, currentTerm=%d", rf.ID, savedCurrentTerm)

	votesReceived := 1
	for _, peer := range rf.peers {
		go func() {
			input := RequestVoteInput{
				Term:        savedCurrentTerm,
				CandidateId: rf.ID,
			}
			voteOutput, err := rf.sendRequestVote(input, peer.ID)
			if err != nil {
				log.Printf("[%d]error sending vote for peer %d: %+v", rf.ID, peer.ID, err)
				return
			}

			rf.mu.Lock()
			defer rf.mu.Unlock()
			if rf.state != Candidate {
				return
			}
			if voteOutput.Term > savedCurrentTerm {
				// become follower
				rf.becomeFollower(voteOutput.Term)
				return
			} else if voteOutput.Term == savedCurrentTerm {
				if voteOutput.VoteGranted {
					votesReceived++
				}
				if votesReceived*2 > len(rf.peers)+1 {
					// Won the election
					log.Printf("[%d]voted to be leader: %d > %d/2", rf.ID, votesReceived, len(rf.peers)+1)
					rf.becomeLeader()
					return
				}
			}
		}()
	}

}

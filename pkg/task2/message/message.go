package message

type Message interface {
	GetSenderId() int
}
type RequestVoteMessage struct {
	SenderID int
	PeerID   int

	Term        int
	CandidateId int
}

func (m *RequestVoteMessage) GetSenderId() int {
	return m.SenderID
}

type RequestVoteMessageReply struct {
	SenderID int
	PeerID   int

	Term        int
	VoteGranted bool
}

func (m *RequestVoteMessageReply) GetSenderId() int {
	return m.SenderID
}

type HeartbeatMessage struct {
	SenderID int
	PeerID   int

	Term     int
	LeaderId int
}

func (m *HeartbeatMessage) GetSenderId() int {
	return m.SenderID
}

type HeartbeatMessageReply struct {
	SenderID int
	PeerID   int

	Term    int
	Success bool
}

func (m *HeartbeatMessageReply) GetSenderId() int {
	return m.SenderID
}

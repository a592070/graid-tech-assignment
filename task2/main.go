package main

import (
	"bufio"
	"fmt"
	"graid-tech-assignment/pkg/task2/quorum_election"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <number_of_members>")
		os.Exit(1)
	}

	numMembers, err := strconv.Atoi(os.Args[1])
	if err != nil || numMembers < 2 {
		fmt.Println("Number of members must be an integer >= 2")
		os.Exit(1)
	}

	raftCluster := make(map[int]*quorum_election.Raft)
	readyChan := make(chan any)

	for i := 0; i < numMembers; i++ {
		raft := quorum_election.NewRaft(
			i,
			readyChan,
		)
		raftCluster[i] = raft
	}
	for i := 0; i < numMembers; i++ {
		raft := raftCluster[i]
		for j := 0; j < numMembers; j++ {
			if i == j {
				continue
			}
			raft.AppendPeer(j,
				raftCluster[j],
			)
		}
	}

	readyChan <- "start"
	defer close(readyChan)

	fmt.Println(`
Offer the following command:
status: Show the state of cluster
kill n: Kill the member
exit: Exit the program
`)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		if command == "exit" {
			for _, raft := range raftCluster {
				raft.Stop()
			}
			break
		} else if command == "status" {
			for _, raft := range raftCluster {
				id, term, state := raft.Report()
				log.Printf("Member %d: state=%s, term=%d\n", id, state, term)
			}
		} else if strings.HasPrefix(command, "kill") {
			parts := strings.Fields(command)
			if len(parts) == 2 {
				if memberID, err := strconv.Atoi(parts[1]); err == nil {
					if raftCluster[memberID] == nil {
						log.Printf("Member %d not found\n", memberID)
						continue
					}
					raftCluster[memberID].Stop()
				}
			}
		}
	}
}

Write a program that demonstrates quorum election. 
The program should have a specified number of members in the quorum 
and start an interactive mode for the quorum election game.

1. Game steps:
   1. Start the quorum with N members.
   2. Elect one of the members as the quorum leader.
   3. Each member sends heartbeat signals to each other to ensure they are alive.
   4. Identify a member that has failed to respond to the heartbeat by voting. 
      1. Remove the failed member from the quorum.
      2. If the failed member was the leader, go back to step ii.
2. Each member should have an ID starting from 0, 1, 2, and so on.
3. The command "kill 0" should make member 0 unresponsive to others.
4. There are multiple quorum mechanisms available, and you can design a better one according to your requirements. (Hint: Consensus Algorithm, or Centralization)
5. Example output:
   ```text
   # launch binary with specified number of member
   ./main 3
   > Starting quorum with 3 members
   > Member 0: Hi
   > Member 1: Hi
   > Member 2: Hi
   > Member 0: I want to be leader
   > Member 2: Accept member 0 to be leader
   > Member 1: I want to be leader
   > Member 1: Accept member 0 to be leader
   > Member 0 voted to be leader: (2 > 3/2)
   > kill 1
   > Member 0: failed heartbeat with Member 1
   > Member 2: failed heartbeat with Member 1
   > Member 1: kick out of quorum: (2 > current/2)
   > kill 2
   > Member 0: failed heartbeat with Member 1
   > Member 0: no response from other users(timeout)
   > Member 2: kick out of quorum: leader decision
   > Quorum failed: (1 > total/2)
   ```

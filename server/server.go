package server

import (
	"context"
	"graft/graft_rpc"
	"graft/models"
	"log"
	"math"
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

// Should create helper function an mu to lock reset and set of timers/tickers
type Server struct {
	TermTimeout     *time.Timer
	ElectionTimeout *time.Timer
	HeartbeatTicker *time.Ticker
}

func sendRequestVoteRpc(host string, input *graft_rpc.RequestVoteInput) (*graft_rpc.RequestVoteOutput, error) {
	var conn *grpc.ClientConn
	conn, _ = grpc.Dial("127.0.0.1:"+host, grpc.WithInsecure())
	defer conn.Close()

	c := graft_rpc.NewRpcClient(conn)

	res, err := c.RequestVote(context.Background(), input)
	return res, err
}

func sendAppendEntriesRpc(host string, input *graft_rpc.AppendEntriesInput) (*graft_rpc.AppendEntriesOutput, error) {
	var conn *grpc.ClientConn
	conn, _ = grpc.Dial("127.0.0.1:"+host, grpc.WithInsecure())
	defer conn.Close()

	c := graft_rpc.NewRpcClient(conn)

	res, err := c.AppendEntries(context.Background(), input)
	return res, err
}

func election(server *Server, state *models.ServerState) {
	// Missed cluster leader heartbeat (leader down, network partition, delays etc...)
	// Start new election
	state.SwitchRole(models.Candidate)          // become candidate
	state.CurrentTerm += 1                      // increment term
	state.PersistentState.VotedFor = state.Name // vote for self

	// Reset election timer
	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(300-150) + 150) * 10
	server.ElectionTimeout = time.NewTimer(time.Duration(timeout) * time.Millisecond)

	// Request votes to each nodes
	totalNodes := len(state.Nodes) + 1 // add self
	quorum := math.Ceil(float64(totalNodes) / 2.0)
	voteGranted := 1 // self

	for _, node := range state.Nodes {
		// Send requestVote RPC
		res, err := sendRequestVoteRpc(node.Host, &graft_rpc.RequestVoteInput{
			Term:         int32(state.CurrentTerm),
			CandidateId:  state.Name,
			LastLogIndex: int32(len(state.PersistentState.Log)),
			LastLogTerm:  int32(state.PersistentState.Log[len(state.PersistentState.Log)-1].Term),
		})
		if err != nil {
			continue
		}

		if res.VoteGranted {
			log.Printf("grant vote by %v \n", node.Host)
			voteGranted += 1
		}
		if res.Term > int32(state.CurrentTerm) {
			// Fallback to follower
			log.Println("fallback to follower, received vote with higher term")
			state.FallbackToFollower(uint16(res.Term))
			server.ElectionTimeout.Stop() // stop election
			break
		}
	}

	// Become leader if quorum reached
	if voteGranted >= int(quorum) && state.IsRole(models.Candidate) {
		log.Printf("become leader, quorum: %v, votes %v \n", quorum, voteGranted)
		state.SwitchRole(models.Leader)
		server.ElectionTimeout.Stop() // stop election
	}

}

func heartbeat(server *Server, state *models.ServerState) {
	// Receive heartbeat from cluster leader
	server.ElectionTimeout.Stop()
	server.TermTimeout.Reset(models.HEARTBEAT * time.Millisecond)
}

func sendHeartbeat(server *Server, state *models.ServerState) {
	for _, node := range state.Nodes {
		// Send heartbeat rpc
		res, err := sendAppendEntriesRpc(node.Host, &graft_rpc.AppendEntriesInput{
			Term:     int32(state.CurrentTerm),
			LeaderId: state.Name,
		})
		if err != nil {
			continue
		}

		if res.Term > int32(state.CurrentTerm) {
			// Fallback to follower
			state.FallbackToFollower(uint16(res.Term))
			server.TermTimeout.Reset(models.HEARTBEAT * time.Millisecond)
			break
		}
	}
}

func (s *Server) Start(state *models.ServerState) {
	for {
		select {
		case <-s.HeartbeatTicker.C:
			if state.IsRole(models.Leader) {
				log.Println("PULSE")
				sendHeartbeat(s, state)
			}

		case <-s.ElectionTimeout.C:
			log.Println("ELECTION_TIMEOUT")
			if state.IsRole(models.Candidate) {
				log.Println("ELECTION_TIMEOUT: restart election")
				election(s, state)
			}

		case <-s.TermTimeout.C:
			log.Println("TERM_TIMEOUT")
			if state.IsRole(models.Follower) {
				log.Println("TERM_TIMEOUT: start election")
				election(s, state)
			}

		case <-state.Heartbeat:
			log.Println("HEARTBEAT")
			heartbeat(s, state)
		}
	}
}

func NewServer() *Server {
	return &Server{
		TermTimeout:     time.NewTimer(models.HEARTBEAT * time.Millisecond),
		ElectionTimeout: time.NewTimer(3000 * time.Millisecond),
		HeartbeatTicker: time.NewTicker(1000 * time.Millisecond),
	}
}

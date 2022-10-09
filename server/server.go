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

type Server struct {
	TermTimeout     *time.Ticker
	ElectionTimeout *time.Timer
}

func sendRequestVoteRpc(host string, input *graft_rpc.RequestVoteInput) (*graft_rpc.RequestVoteOutput, error) {
	var conn *grpc.ClientConn
	conn, _ = grpc.Dial("127.0.0.1:"+host, grpc.WithInsecure())
	defer conn.Close()

	c := graft_rpc.NewRpcClient(conn)

	res, err := c.RequestVote(context.Background(), input)
	return res, err
}

func election(server *Server, state *models.ServerState) {
	// Missed cluster leader heartbeat (leader down, network partition, delays etc...)
	// Start new election
	state.Role = models.Candidate               // become candidate
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
			state.Role = models.Follower
			server.ElectionTimeout.Stop() // stop election
			break
		}
	}

	// Become leader if quorum reached
	if voteGranted >= int(quorum) && state.Role == models.Candidate {
		log.Printf("become leader, quorum: %v, votes %v \n", quorum, voteGranted)
		state.Role = models.Leader
		server.ElectionTimeout.Stop() // stop election
	}

}

func heartbeat(server *Server, state *models.ServerState) {
	// Receive heartbeat from cluster leader
	server.ElectionTimeout.Stop()
	server.TermTimeout.Stop()
	server.TermTimeout = time.NewTicker(models.HEARTBEAT * time.Millisecond)
}

func (s *Server) Start(state *models.ServerState) {
	for {
		select {
		case <-s.ElectionTimeout.C:
			log.Println("election timeout, start election")
			s.TermTimeout.Stop() // prevent two elections ocurring at the same time
			election(s, state)
			s.TermTimeout = time.NewTicker(models.HEARTBEAT * time.Millisecond)

		case <-s.TermTimeout.C:
			if state.Role == models.Follower {
				log.Println("term timeout, start election")
				election(s, state)
			}

		case <-state.Heartbeat:
			log.Println("heartbeat")
			heartbeat(s, state)
		}
	}
}

// TODO: leader should send heartbeat request to followers

func NewServer() *Server {
	return &Server{
		TermTimeout:     time.NewTicker(models.HEARTBEAT * time.Millisecond),
		ElectionTimeout: time.NewTimer(3000 * time.Millisecond),
	}
}

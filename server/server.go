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
	// timeout := rand.Intn(300-150) + 150
	timeout := 1000
	server.ElectionTimeout = time.NewTimer(time.Duration(timeout) * time.Millisecond)

	// Request votes to each nodes
	totalNodes := len(state.Nodes) + 1 // add self
	quorum := math.Ceil(float64(totalNodes) / 2.0)
	voteGranted := 1 // self
	for _, node := range state.Nodes {
		// Send requestVote RPC
		res, err := sendRequestVoteRpc(node.Host, &graft_rpc.RequestVoteInput{})
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
			break
		}
	}

	// Become leader if quorum reached
	if voteGranted >= int(quorum) && state.Role == models.Candidate {
		log.Printf("become leader, quorum: %v, votes %v \n", quorum, voteGranted)
		state.Role = models.Leader
	}
}

func heartbeat(server *Server, state *models.ServerState) {
	// Receive heartbeat from cluster leader
	server.TermTimeout.Stop()
	server.TermTimeout = time.NewTicker(models.HEARTBEAT * time.Millisecond)
}

func (s *Server) Start(state *models.ServerState) {
	for {
		select {
		case <-s.ElectionTimeout.C:
		case <-s.TermTimeout.C:
			log.Println("timeout, start election")
			election(s, state)

		case <-state.Heartbeat:
			log.Println("heartbeat")
			heartbeat(s, state)
		}
	}
}

func NewServer() *Server {
	return &Server{
		TermTimeout:     time.NewTicker(models.HEARTBEAT * time.Millisecond),
		ElectionTimeout: time.NewTimer(1000 * time.Millisecond),
	}
}

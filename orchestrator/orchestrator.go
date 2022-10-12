package orchestrator

import (
	"graft/api"
	"graft/api/graft_rpc"
	"graft/models"
	"log"
	"math/rand"
	"sync"
	"time"
)

const HEARTBEAT_TICKER = 50  // ms
const ELECTION_TIMEOUT = 350 // ms

type EventOrchestrator struct {
	serverState *models.ServerState

	heartbeatTicker *time.Ticker
	electionTimer   *time.Timer
	mu              sync.Mutex
}

func NewEventOrchestrator(state *models.ServerState) *EventOrchestrator {
	return &EventOrchestrator{
		serverState: state,

		heartbeatTicker: time.NewTicker(HEARTBEAT_TICKER * time.Millisecond),
		electionTimer:   time.NewTimer(ELECTION_TIMEOUT * time.Millisecond),
	}
}

func (och *EventOrchestrator) Start() {
	log.Println("START EVENT ORCHESTRATOR")

	counter := 0
	for {
		select {
		case <-och.heartbeatTicker.C:
			if och.serverState.IsLeader() {
				och.sendHeartbeat()
			}

		case <-och.electionTimer.C:
			if och.serverState.IsFollower() || och.serverState.IsCandidate() {
				log.Println("TIMEOUT")
				och.startElection()
			}

		case <-och.serverState.Heartbeat():
			counter += 1
			if counter%(HEARTBEAT_TICKER*3) == 0 {
				log.Println("HEARTBEAT")
			}
			och.resetElectionTimeout()
		}
	}
}

func (och *EventOrchestrator) resetElectionTimeout() {
	och.mu.Lock()
	defer och.mu.Unlock()

	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(ELECTION_TIMEOUT/2) + ELECTION_TIMEOUT/2)
	och.electionTimer.Reset(time.Duration(timeout) * time.Millisecond)
}

func (och *EventOrchestrator) startElection() {
	state := och.serverState
	state.RaiseToCandidate()
	och.resetElectionTimeout()

	votesGranted := 1 // vote for self
	voteInput := &graft_rpc.RequestVoteInput{
		CandidateId:  string(state.Name),
		Term:         int32(state.CurrentTerm),
		LastLogIndex: int32(state.LastLogIndex()),
		LastLogTerm:  int32(state.LastLogTerm()),
	}

	var m sync.Mutex
	var wg sync.WaitGroup
	for _, node := range state.Nodes {
		wg.Add(1)
		go (func(host string, w *sync.WaitGroup) {
			defer w.Done()
			if res, err := api.SendRequestVoteRpc(host, voteInput); err == nil {
				if res.Term > int32(state.CurrentTerm) {
					state.DowngradeToFollower(uint16(res.Term))
					return
				}
				if res.VoteGranted {
					m.Lock()
					defer m.Unlock()
					votesGranted += 1
				}
			}
		})(node.Host, &wg)
	}
	wg.Wait()

	if votesGranted >= state.Quorum() && state.IsCandidate() {
		state.PromoteToLeader()
	}
}

func (och *EventOrchestrator) sendHeartbeat() {
	state := och.serverState

	heartbeatInput := &graft_rpc.AppendEntriesInput{
		Term:     int32(state.CurrentTerm),
		LeaderId: state.Name,
	}

	var wg sync.WaitGroup
	for _, node := range state.Nodes {
		wg.Add(1)
		go (func(host string, w *sync.WaitGroup) {
			defer w.Done()
			if res, err := api.SendAppendEntriesRpc(host, heartbeatInput); err == nil {
				if res.Term > int32(state.CurrentTerm) {
					state.DowngradeToFollower(uint16(res.Term))
					return
				}
			}
		})(node.Host, &wg)
	}
	wg.Wait()
}

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

const HEARTBEAT_FREQ = 50
const TERM_TIMEOUT = 250
const ELECTION_TIMEOUT = 350

type EventOrchestrator struct {
	heartbeatTicker *time.Ticker
	termTicker      *time.Ticker
	electionTimer   *time.Timer
	mu              sync.Mutex
}

func NewEventOrchestrator() *EventOrchestrator {
	return &EventOrchestrator{
		heartbeatTicker: time.NewTicker(HEARTBEAT_FREQ * time.Millisecond),
		termTicker:      time.NewTicker(TERM_TIMEOUT * time.Millisecond),
		electionTimer:   time.NewTimer(ELECTION_TIMEOUT * time.Millisecond),
	}
}

func (och *EventOrchestrator) start(state *models.ServerState) {
	counter := 0
	for {
		select {
		case <-och.heartbeatTicker.C:
			if state.IsLeader() {
				och.sendHeartbeat(state)
			}

		case <-och.electionTimer.C:
			if state.IsCandidate() {
				log.Println("ELECTION_TIMEOUT")
				och.startElection(state)
			}

		case <-och.termTicker.C:
			if state.IsFollower() {
				log.Println("TERM_TIMEOUT")
				och.startElection(state)
			}

		case <-state.Heartbeat():
			counter += 1
			if counter%(HEARTBEAT_FREQ*3) == 0 {
				log.Println("HEARTBEAT")
			}
			och.heartbeat(state)
		}
	}
}

func (och *EventOrchestrator) heartbeat(state *models.ServerState) {
	och.mu.Lock()
	defer och.mu.Unlock()

	och.termTicker.Reset(TERM_TIMEOUT * time.Millisecond)
}

func (och *EventOrchestrator) resetElectionTimeout() {
	och.mu.Lock()
	defer och.mu.Unlock()

	rand.Seed(time.Now().UnixNano())
	timeout := (rand.Intn(ELECTION_TIMEOUT/2) + ELECTION_TIMEOUT/2)
	och.electionTimer.Reset(time.Duration(timeout) * time.Millisecond)
}

func (och *EventOrchestrator) startElection(state *models.ServerState) {
	state.RaiseToCandidate()
	och.resetElectionTimeout()

	votesGranted := 1 // vote for self
	voteInput := &graft_rpc.RequestVoteInput{
		Term:         int32(state.CurrentTerm),
		CandidateId:  string(state.Name),
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

func (och *EventOrchestrator) sendHeartbeat(state *models.ServerState) {
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
				}
			}
		})(node.Host, &wg)
	}
	wg.Wait()
}

func StartEventOrchestrator(state *models.ServerState) *EventOrchestrator {
	log.Println("START EVENT ORCHESTRATOR")
	orchestrator := NewEventOrchestrator()
	go orchestrator.start(state)
	return orchestrator
}

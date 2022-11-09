package runner

import (
	"sync"
	"sync/atomic"

	"graft/pkg/domain/entity"
	domain "graft/pkg/domain/service"

	log "github.com/sirupsen/logrus"
)

type synchronise struct {
	save      sync.Mutex
	commit    sync.Mutex
	election  sync.Mutex
	heartbeat sync.Mutex
}

type service struct {
	clusterNode *domain.ClusterNode
	timeout     *entity.Timeout
	repo        repository
	persist     persister
	sync        synchronise
}

func NewService(
	clusterNode *domain.ClusterNode,
	timeout *entity.Timeout,
	repository repository,
	persister persister,
) *service {
	return &service{
		repo:        repository,
		timeout:     timeout,
		clusterNode: clusterNode,
		persist:     persister,
	}
}

func (s *service) Run() {
	select {
	case role := <-s.clusterNode.ShiftRole:
		go s.runAs(role)

	case <-s.clusterNode.SaveState:
		go s.saveState()

	case <-s.clusterNode.Commit:
		go s.commit()

	case <-s.clusterNode.ResetElectionTimer:
		go s.timeout.ResetElectionTimer()

	case <-s.clusterNode.ResetLeaderTicker:
		go s.timeout.ResetLeaderTicker()
	}
}

func (s *service) runAs(role entity.Role) {
	switch role {
	case entity.Follower:
		s.runFollower(s.clusterNode)
	case entity.Candidate:
		s.runCandidate(s.clusterNode)
	case entity.Leader:
		s.runLeader(s.clusterNode)
	}
}

func (s *service) commit() {
	s.sync.commit.Lock()
	defer s.sync.commit.Unlock()
	s.clusterNode.ApplyLogs()
}

func (s *service) saveState() {
	s.sync.save.Lock()
	defer s.sync.save.Unlock()
	s.persist.Save(
		s.clusterNode.CurrentTerm(),
		s.clusterNode.VotedFor(),
		s.clusterNode.MachineLogs(),
	)
}

func (s *service) runFollower(f follower) {
	for range s.timeout.ElectionTimer.C {
		f.UpgradeCandidate()
		return
	}
}

func (s *service) runCandidate(c candidate) {
	if wonElection := s.runElection(c); wonElection {
		c.UpgradeLeader()
	} else {
		// Restart election until: upgrade / downgrade
		for range s.timeout.ElectionTimer.C {
			if wonElection := s.runElection(c); wonElection {
				c.UpgradeLeader()
				return
			}
		}
	}
}

func (s *service) runLeader(l leader) {
	for range s.timeout.LeaderTicker.C {
		if quorumReached := s.synchronizeLogs(l); !quorumReached {
			log.Debug("LEADER STEP DOWN")
			l.DowngradeFollower(l.GetState().CurrentTerm())
			return
		}
	}
}

func (s *service) runElection(c candidate) bool {
	s.sync.election.Lock()
	defer s.sync.election.Unlock()

	c.IncrementCandidateTerm()
	state := c.GetState()
	input := c.RequestVoteInput()
	quorum := c.Quorum()

	var prevotesGranted uint32 = 1 // vote for self
	preVoteRoutine := func(p entity.Peer) {
		if res, err := s.repo.PreVote(p, &input); err == nil {
			if res.Term > state.CurrentTerm() {
				c.DowngradeFollower(res.Term)
				return
			}
			if res.VoteGranted {
				atomic.AddUint32(&prevotesGranted, 1)
			}
		}
	}
	c.Broadcast(preVoteRoutine)

	if int(prevotesGranted) < quorum {
		return false
	}
	log.Debug("GATHER VOTES")

	var votesGranted uint32 = 1 // vote for self
	gatherVotesRoutine := func(p entity.Peer) {
		if res, err := s.repo.RequestVote(p, &input); err == nil {
			if res.Term > state.CurrentTerm() {
				c.DowngradeFollower(res.Term)
				return
			}
			if res.VoteGranted {
				atomic.AddUint32(&votesGranted, 1)
			}
		}
	}
	c.Broadcast(gatherVotesRoutine)

	return int(votesGranted) >= quorum
}

func (s *service) synchronizeLogs(l leader) bool {
	s.sync.heartbeat.Lock()
	defer s.sync.heartbeat.Unlock()

	state := l.GetState()
	quorum := l.Quorum()

	var peersAlive uint32 = 1 // self
	synchroniseLogsRoutine := func(p entity.Peer) {
		input := l.AppendEntriesInput(p.Id)

		if res, err := s.repo.AppendEntries(p, &input); err == nil {
			atomic.AddUint32(&peersAlive, 1)

			if res.Term > state.CurrentTerm() {
				l.DowngradeFollower(res.Term)
				return
			}
			if res.Success {
				leaderLastLogIndex := state.LastLogIndex()
				if state.ShouldUpdatePeerIndex(p.Id) {
					l.SetNextMatchIndex(p.Id, leaderLastLogIndex)
				}
			} else {
				l.DecrementNextIndex(p.Id)
			}
		}
	}
	l.Broadcast(synchroniseLogsRoutine)

	newCommitIndex := l.ComputeNewCommitIndex()
	l.SetCommitIndex(newCommitIndex)
	return int(peersAlive) >= quorum
}

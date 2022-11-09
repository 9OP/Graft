package runner

import (
	"sync/atomic"

	"graft/pkg/domain/entity"
	domain "graft/pkg/domain/service"

	log "github.com/sirupsen/logrus"
)

type synchronise struct {
	running    uint32
	saving     uint32
	committing uint32
}

func (s *synchronise) CanRun() bool {
	return atomic.CompareAndSwapUint32(&s.running, 0, 1)
}

func (s *synchronise) CanSave() bool {
	return atomic.CompareAndSwapUint32(&s.saving, 0, 1)
}

func (s *synchronise) CanCommit() bool {
	return atomic.CompareAndSwapUint32(&s.committing, 0, 1)
}

func (s *synchronise) ReleaseRun() {
	atomic.StoreUint32(&s.running, 0)
}

func (s *synchronise) ReleaseSave() {
	atomic.StoreUint32(&s.saving, 0)
}

func (s *synchronise) ReleaseCommit() {
	atomic.StoreUint32(&s.committing, 0)
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

func (s *service) Dispatch() {
	select {
	case role := <-s.clusterNode.ShiftRole:
		if s.sync.CanRun() {
			go s.runAs(role)
		}

	case <-s.clusterNode.SaveState:
		if s.sync.CanSave() {
			go s.saveState()
		}

	case <-s.clusterNode.Commit:
		if s.sync.CanCommit() {
			go s.commit()
		}

	case <-s.clusterNode.ResetElectionTimer:
		go s.timeout.ResetElectionTimer()

	case <-s.clusterNode.ResetLeaderTicker:
		go s.timeout.ResetLeaderTicker()
	}
}

func (s *service) runAs(role entity.Role) {
	defer s.sync.ReleaseRun()

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
	defer s.sync.ReleaseCommit()
	s.clusterNode.ApplyLogs()
}

func (s *service) saveState() {
	defer s.sync.ReleaseSave()
	s.persist.Save(
		s.clusterNode.CurrentTerm(),
		s.clusterNode.VotedFor(),
		s.clusterNode.MachineLogs(),
	)
}

func (s *service) runFollower(f follower) {
	for {
		select {
		case <-s.timeout.ElectionTimer.C:
			f.UpgradeCandidate()
			return
		case <-s.clusterNode.ShiftRole:
			// quit
			log.Debug("QUIT runFollower")
			return
		}
	}
}

func (s *service) runCandidate(c candidate) {
	if wonElection := s.runElection(c); wonElection {
		c.UpgradeLeader()
		return
	}

	// Restart election
	for {
		select {
		case <-s.timeout.ElectionTimer.C:
			if wonElection := s.runElection(c); wonElection {
				c.UpgradeLeader()
				return
			}
		case <-s.clusterNode.ShiftRole:
			// quit
			log.Debug("QUIT runCandidate")
			return
		}
	}
}

func (s *service) runLeader(l leader) {
	for {
		select {
		case <-s.timeout.LeaderTicker.C:
			if quorumReached := s.synchronizeLogs(l); !quorumReached {
				log.Debug("LEADER STEP DOWN")
				l.DowngradeFollower(l.GetState().CurrentTerm())
				return
			}
		case <-s.clusterNode.ShiftRole:
			// quit
			log.Debug("QUIT runLeader")
			return
		}
	}
}

func (s *service) runElection(c candidate) bool {
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

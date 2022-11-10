package runner

import (
	"sync/atomic"

	"graft/pkg/domain/entity"
	domain "graft/pkg/domain/service"

	log "github.com/sirupsen/logrus"
)

type binaryLock struct {
	c uint32
}

func (l *binaryLock) Lock() bool {
	return atomic.CompareAndSwapUint32(&l.c, 0, 1)
}

func (l *binaryLock) Unlock() {
	atomic.StoreUint32(&l.c, 0)
}

type synchronise struct {
	running    binaryLock
	saving     binaryLock
	committing binaryLock
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
	for {
		switch s.clusterNode.Role() {
		case entity.Follower:
			s.runFollower()
		case entity.Candidate:
			s.runCandidate()
		case entity.Leader:
			s.runLeader()
		}
	}
}

func (s *service) followerFlow() {
	// Prevent multiple flow from running together
	if s.sync.running.Lock() {
		defer s.sync.running.Unlock()

		s.clusterNode.UpgradeCandidate()
		if wonElection := s.runElection(); wonElection {
			s.clusterNode.UpgradeLeader()
		}
	}
}

func (s *service) candidateFlow() {
	// Prevent multiple flow from running together
	if s.sync.running.Lock() {
		defer s.sync.running.Unlock()

		if wonElection := s.runElection(); wonElection {
			s.clusterNode.UpgradeLeader()
		}
	}
}

func (s *service) leaderFlow() {
	// Prevent multiple flow from running together
	if s.sync.running.Lock() {
		defer s.sync.running.Unlock()

		if quorumReached := s.heartbeat(); !quorumReached {
			log.Debug("STEP DOWN")
			s.clusterNode.DowngradeFollower(s.clusterNode.CurrentTerm())
		}
	}
}

func (s *service) runFollower() {
	for s.clusterNode.Role() == entity.Follower {
		select {
		case <-s.timeout.ElectionTimer.C:
			go s.followerFlow()

		case <-s.clusterNode.ResetElectionTimer:
			go s.timeout.ResetElectionTimer()

		case <-s.clusterNode.SaveState:
			go s.saveState()

		case <-s.clusterNode.ShiftRole:
			return
		}
	}
}

func (s *service) runCandidate() {
	for s.clusterNode.Role() == entity.Candidate {
		select {
		case <-s.timeout.ElectionTimer.C:
			go s.candidateFlow()

		case <-s.clusterNode.ResetElectionTimer:
			go s.timeout.ResetElectionTimer()

		case <-s.clusterNode.SaveState:
			go s.saveState()

		case <-s.clusterNode.ShiftRole:
			return
		}
	}
}

func (s *service) runLeader() {
	for s.clusterNode.Role() == entity.Leader {
		select {
		case <-s.timeout.LeaderTicker.C:
			go s.leaderFlow()

		// case <-s.clusterNode.SynchronizeLogs:
		// 	go s.timeout.ResetLeaderTicker()
		// 	go s.leaderFlow()

		case <-s.clusterNode.Commit:
			go s.commit()

		case <-s.clusterNode.SaveState:
			go s.saveState()

		case <-s.clusterNode.ResetLeaderTicker:
			go s.timeout.ResetLeaderTicker()

		case <-s.clusterNode.ShiftRole:
			return
		}
	}
}

func (s *service) commit() {
	if s.sync.committing.Lock() {
		defer s.sync.committing.Unlock()

		s.clusterNode.ApplyLogs()
	}
}

func (s *service) saveState() {
	if s.sync.saving.Lock() {
		defer s.sync.saving.Unlock()

		s.persist.Save(
			s.clusterNode.CurrentTerm(),
			s.clusterNode.VotedFor(),
			s.clusterNode.MachineLogs(),
		)
	}
}

func (s *service) runElection() bool {
	s.clusterNode.IncrementCandidateTerm()

	if !s.preVote() {
		return false
	}

	return s.requestVote()
}

func (s *service) preVote() bool {
	input := s.clusterNode.RequestVoteInput()
	quorum := s.clusterNode.Quorum()
	var prevotesGranted uint32 = 1 // vote for self

	preVoteRoutine := func(p entity.Peer) {
		if res, err := s.repo.PreVote(p, &input); err == nil {
			if res.Term > s.clusterNode.CurrentTerm() {
				s.clusterNode.DowngradeFollower(res.Term)
				return
			}
			if res.VoteGranted {
				atomic.AddUint32(&prevotesGranted, 1)
			}
		}
	}
	s.clusterNode.Broadcast(preVoteRoutine)

	isCandidate := s.clusterNode.Role() == entity.Candidate
	quorumReached := int(prevotesGranted) >= quorum
	return isCandidate && quorumReached
}

func (s *service) requestVote() bool {
	input := s.clusterNode.RequestVoteInput()
	quorum := s.clusterNode.Quorum()
	var votesGranted uint32 = 1 // vote for self

	gatherVotesRoutine := func(p entity.Peer) {
		if res, err := s.repo.RequestVote(p, &input); err == nil {
			if res.Term > s.clusterNode.CurrentTerm() {
				s.clusterNode.DowngradeFollower(res.Term)
				return
			}
			if res.VoteGranted {
				atomic.AddUint32(&votesGranted, 1)
			}
		}
	}
	s.clusterNode.Broadcast(gatherVotesRoutine)

	isCandidate := s.clusterNode.Role() == entity.Candidate
	quorumReached := int(votesGranted) >= quorum
	return isCandidate && quorumReached
}

func (s *service) heartbeat() bool {
	quorumReached := s.synchronizeLogs()

	if s.clusterNode.Role() == entity.Leader {
		s.clusterNode.SetCommitIndex(
			s.clusterNode.ComputeNewCommitIndex(),
		)
	}

	return quorumReached
}

func (s *service) synchronizeLogs() bool {
	quorum := s.clusterNode.Quorum()
	var peersAlive uint32 = 1 // self

	synchroniseLogsRoutine := func(p entity.Peer) {
		input := s.clusterNode.AppendEntriesInput(p.Id)
		if res, err := s.repo.AppendEntries(p, &input); err == nil {
			atomic.AddUint32(&peersAlive, 1)
			if res.Term > s.clusterNode.CurrentTerm() {
				s.clusterNode.DowngradeFollower(res.Term)
				return
			}
			if res.Success {
				leaderLastLogIndex := s.clusterNode.LastLogIndex()
				if s.clusterNode.ShouldUpdatePeerIndex(p.Id) {
					s.clusterNode.SetNextMatchIndex(p.Id, leaderLastLogIndex)
				}
			} else {
				s.clusterNode.DecrementNextIndex(p.Id)
			}
		}
	}
	s.clusterNode.Broadcast(synchroniseLogsRoutine)

	quorumReached := int(peersAlive) >= quorum
	return quorumReached
}

package runner

import (
	"sync/atomic"

	"graft/pkg/domain"

	log "github.com/sirupsen/logrus"
)

type synchronise struct {
	running    binaryLock
	saving     binaryLock
	committing binaryLock
}

type service struct {
	clusterNode *domain.Node
	timeout     *timeout
	repo        repository
	persist     persister
	sync        synchronise
}

func NewService(
	clusterNode *domain.Node,
	repository repository,
	persister persister,
	electionDuration int,
	leaderHearbeat int,
) *service {
	timeout := newTimeout(electionDuration, leaderHearbeat)
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
		case domain.Follower:
			s.runFollower()
		case domain.Candidate:
			s.runCandidate()
		case domain.Leader:
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
	for s.clusterNode.Role() == domain.Follower {
		select {
		case <-s.timeout.ElectionTimer.C:
			go s.followerFlow()

		case <-s.clusterNode.Commit:
			go s.commit()

		case <-s.clusterNode.ResetElectionTimer:
			go s.timeout.resetElectionTimer()

		case <-s.clusterNode.SaveState:
			go s.saveState()

		case <-s.clusterNode.ShiftRole:
			return
		}
	}
}

func (s *service) runCandidate() {
	for s.clusterNode.Role() == domain.Candidate {
		select {
		case <-s.timeout.ElectionTimer.C:
			go s.candidateFlow()

		case <-s.clusterNode.ResetElectionTimer:
			go s.timeout.resetElectionTimer()

		case <-s.clusterNode.Commit:
			go s.commit()

		case <-s.clusterNode.SaveState:
			go s.saveState()

		case <-s.clusterNode.ShiftRole:
			return
		}
	}
}

func (s *service) runLeader() {
	for s.clusterNode.Role() == domain.Leader {
		select {
		case <-s.timeout.LeaderTicker.C:
			go s.leaderFlow()

		// Should batch execution
		case <-s.clusterNode.SynchronizeLogs:
			go (func() {
				s.timeout.resetLeaderTicker()
				s.leaderFlow()
			})()

		case <-s.clusterNode.Commit:
			go s.commit()

		case <-s.clusterNode.SaveState:
			go s.saveState()

		case <-s.clusterNode.ResetLeaderTicker:
			go s.timeout.resetLeaderTicker()

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
			s.clusterNode.Logs(),
		)
	}
}

func (s *service) runElection() bool {
	// Prevent candidate with outdated logs
	// From running election
	if !s.preVote() {
		log.Debug("FAILED PRE-VOTE")
		return false
	}

	s.clusterNode.IncrementCandidateTerm()

	return s.requestVote()
}

func (s *service) preVote() bool {
	input := s.clusterNode.RequestVoteInput()
	quorum := s.clusterNode.Quorum()
	var prevotesGranted uint32 = 1 // vote for self

	preVoteRoutine := func(p domain.Peer) {
		if res, err := s.repo.PreVote(p, &input); err == nil {
			if res.VoteGranted {
				atomic.AddUint32(&prevotesGranted, 1)
			}
		}
	}
	s.clusterNode.Broadcast(preVoteRoutine)

	quorumReached := int(prevotesGranted) >= quorum
	return quorumReached
}

func (s *service) requestVote() bool {
	input := s.clusterNode.RequestVoteInput()
	quorum := s.clusterNode.Quorum()
	var votesGranted uint32 = 1 // vote for self

	gatherVotesRoutine := func(p domain.Peer) {
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

	isCandidate := s.clusterNode.Role() == domain.Candidate
	quorumReached := int(votesGranted) >= quorum
	return isCandidate && quorumReached
}

func (s *service) heartbeat() bool {
	quorumReached := s.synchronizeLogs()

	if s.clusterNode.Role() == domain.Leader {
		s.clusterNode.SetCommitIndex(
			s.clusterNode.ComputeNewCommitIndex(),
		)
	}

	return quorumReached
}

func (s *service) synchronizeLogs() bool {
	quorum := s.clusterNode.Quorum()
	var peersAlive uint32 = 1 // self

	synchroniseLogsRoutine := func(p domain.Peer) {
		input := s.clusterNode.AppendEntriesInput(p.Id)
		if res, err := s.repo.AppendEntries(p, &input); err == nil {
			atomic.AddUint32(&peersAlive, 1)
			if res.Term > s.clusterNode.CurrentTerm() {
				s.clusterNode.DowngradeFollower(res.Term)
				return
			}
			if res.Success {
				// newPeerLastLogIndex is always the sum of len(entries) and prevLogIndex
				// Even if peer was send logs it already has, because prevLogIndex is the number
				// of logs already contained at least in peer, and len(entries) is the additional
				// entries accepted.
				newPeerLastLogIndex := input.PrevLogIndex + uint32(len(input.Entries))
				s.clusterNode.SetNextMatchIndex(p.Id, newPeerLastLogIndex)
			} else {
				s.clusterNode.DecrementNextIndex(p.Id)
			}
		}
	}
	s.clusterNode.Broadcast(synchroniseLogsRoutine)

	quorumReached := int(peersAlive) >= quorum
	return quorumReached
}

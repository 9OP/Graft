package runner

import "fmt"

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) RunFollower(follower Follower) {
run:
	for {
		select {
		case <-follower.Timeout():
			follower.UpgradeCandidate()
			break run
		default:
			continue
		}
	}

}

func (s *Service) RunCandidate(candidate Candidate) {
	fmt.Println("RUN AS CANDIDATE")
	// Start election
	// if timeout, restart election
	// if quorum reached, upgradeLeader
	//candidate.Broadcast()
}
func (s *Service) RunLeader(leader Leader) {}

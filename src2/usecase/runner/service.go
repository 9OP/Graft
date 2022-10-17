package runner

import "fmt"

type Service struct{}

func NewService() *Service {
	return &Service{}
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
}
func (s *Service) RunLeader(leader Leader) {}

package entity

type AppendEntriesInput struct {
	LeaderId     string
	Entries      []string
	Term         uint32
	PrevLogIndex uint32
	PrevLogTerm  uint32
	LeaderCommit uint32
}
type AppendEntriesOutput struct {
	Term    uint32
	Success bool
}

type RequestVoteInput struct {
	CandidateId  string
	Term         uint32
	LastLogIndex uint32
	LastLogTerm  uint32
}
type RequestVoteOutput struct {
	Term        uint32
	VoteGranted bool
}

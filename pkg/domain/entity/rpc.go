package entity

type EvalResult struct {
	Out []byte
	Err error
}

type LogEntry struct {
	Term  uint32          `json:"term"`
	Value string          `json:"value"`
	C     chan EvalResult `json:"-"`
}

type AppendEntriesInput struct {
	LeaderId     string
	Entries      []LogEntry
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

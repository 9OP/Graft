package domain

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

type LogEntry struct {
	Term  uint32          `json:"term"`
	Value string          `json:"value"`
	Type  string          `json:"type"`
	C     chan EvalResult `json:"-"`
}

func (l LogEntry) Copy() LogEntry {
	return LogEntry{
		Term:  l.Term,
		Value: l.Value,
		Type:  l.Type,
		C:     l.C,
	}
}

type EvalResult struct {
	Out []byte
	Err error
}

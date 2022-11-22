package domain

import "fmt"

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

// TODO add Date value
type LogEntry struct {
	Index uint64          `json:"index"`
	Term  uint32          `json:"term"`
	Data  []byte          `json:"value"`
	Type  LogType         `json:"type"`
	C     chan EvalResult `json:"-"`
}

type EvalResult struct {
	Out []byte
	Err error
}

type LogType uint8

const (
	LogCommand LogType = iota
	LogQuery
	LogNoop
	LogConfiguration
)

func (lt LogType) String() string {
	switch lt {
	case LogCommand:
		return "LogCommand"
	case LogQuery:
		return "LogQuery"
	case LogNoop:
		return "LogNoop"
	case LogConfiguration:
		return "LogConfiguration"
	default:
		return fmt.Sprintf("%d", lt)
	}
}

type ApiCommand struct {
	Type LogType
	Data []byte
}

type ConfigurationType uint8

const (
	ConfAddPeer ConfigurationType = iota
	ConfActivatePeer
	ConfDeactivatePeer
	ConfRemovePeer
)

type ConfigurationUpdate struct {
	Type ConfigurationType
	Peer
}

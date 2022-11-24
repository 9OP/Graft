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

type ExecuteInput struct {
	Type LogType
	Data []byte
}

type ExecuteOutput struct {
	Out []byte
	Err error
}

// TODO add Date value
type LogEntry struct {
	Index uint64             `json:"index"`
	Term  uint32             `json:"term"`
	Data  []byte             `json:"value"`
	Type  LogType            `json:"type"`
	C     chan ExecuteOutput `json:"-"`
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

type ConfigurationUpdateType uint8

const (
	ConfAddPeer ConfigurationUpdateType = iota
	ConfActivatePeer
	ConfDeactivatePeer
	ConfRemovePeer
)

type ConfigurationUpdate struct {
	Type ConfigurationUpdateType
	Peer
}

type ClusterConfiguration struct {
	Peers           Peers
	LeaderId        string
	ElectionTimeout int
	LeaderHeartbeat int
}

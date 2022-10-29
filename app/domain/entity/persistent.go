package entity

import "sync"

type Persistent struct {
	CurrentTerm uint32       `json:"current_term"`
	VotedFor    string       `json:"voted_for"`
	MachineLogs []MachineLog `json:"machine_logs"`
	mu          sync.RWMutex
}

type MachineLog struct {
	Term  uint32 `json:"term"`
	Value string `json:"value"`
}

func NewPersistent() *Persistent {
	return &Persistent{
		CurrentTerm: 0,
		VotedFor:    "",
		MachineLogs: []MachineLog{},
	}
}

func (p *Persistent) GetCopy() *Persistent {
	p.mu.RLock()
	defer p.mu.RUnlock()

	machineLogs := make([]MachineLog, len(p.MachineLogs))
	copy(machineLogs, p.MachineLogs)

	return &Persistent{
		CurrentTerm: p.CurrentTerm,
		VotedFor:    p.VotedFor,
		MachineLogs: machineLogs,
	}
}

func (p *Persistent) SetVotedFor(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.VotedFor = id
}

func (p *Persistent) SetCurrentTerm(term uint32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.CurrentTerm = term
}

// Delete logs from given index (include)
func (p *Persistent) DeleteLogsFrom(index uint32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	lastLogIndex := p.GetLastLogIndex()
	if index <= uint32(lastLogIndex) {
		p.MachineLogs = p.MachineLogs[:index]
	}
}

func (p *Persistent) AppendLogs(entries []string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	logs := make([]MachineLog, 0, len(entries)+len(p.MachineLogs))
	for _, entry := range entries {
		logs = append(logs, MachineLog{Term: p.CurrentTerm, Value: entry})
	}
	p.MachineLogs = logs
}

// Index starts at 1
func (p *Persistent) GetLogByIndex(index uint32) MachineLog {
	if index <= p.GetLastLogIndex() && index >= 1 {
		return p.MachineLogs[index-1]
	}
	return MachineLog{Term: 0}
}

// Returns index of last log
// Index starts at 1, return value 0 means there are no logs
func (p *Persistent) GetLastLogIndex() uint32 {
	return uint32(len(p.MachineLogs))
}

// Returns lastLogIndex and lastLog
func (p *Persistent) GetLastLog() (uint32, MachineLog) {
	lastLogIndex := p.GetLastLogIndex()
	lastLog := p.GetLogByIndex(lastLogIndex)
	return lastLogIndex, lastLog
}

func (p *Persistent) GetLastLogTerm() uint32 {
	_, lastLog := p.GetLastLog()
	return lastLog.Term
}

// Returns the entire logs stash from given index
func (p *Persistent) GetLogsFromIndex(index uint32) []MachineLog {
	logs := make([]MachineLog, 0, len(p.MachineLogs))
	lastLogIndex := p.GetLastLogIndex()
	for idx := index; idx <= lastLogIndex; idx++ {
		logs = append(logs, p.GetLogByIndex(idx))
	}
	return logs
}

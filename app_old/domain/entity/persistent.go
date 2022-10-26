package entity

import "sync"

type Persistent struct {
	CurrentTerm uint32       `json:"current_term"`
	VotedFor    string       `json:"voted_for"`
	MachineLogs []MachineLog `json:"machine_logs"`
	mu          sync.Mutex
}

type MachineLog struct {
	Term  uint32 `json:"term"`
	Value string `json:"value"`
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

func (p *Persistent) DeleteLogsFrom(index uint32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	lastLogIndex := int32(len(p.MachineLogs))
	if index < uint32(lastLogIndex) {
		p.MachineLogs = p.MachineLogs[:index]
	}
}

func (p *Persistent) AppendLogs(entries []string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	size := len(entries) + len(p.MachineLogs)
	logs := make([]MachineLog, size)
	currentTerm := p.CurrentTerm

	for _, entry := range entries {
		newLog := MachineLog{Term: currentTerm, Value: entry}
		logs = append(logs, newLog)
	}
	p.MachineLogs = logs
}

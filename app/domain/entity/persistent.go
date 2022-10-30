package entity

import (
	"log"
	"sync"
)

type Persistent struct {
	CurrentTerm uint32     `json:"current_term"`
	VotedFor    string     `json:"voted_for"`
	MachineLogs []LogEntry `json:"machine_logs"`
	mu          sync.RWMutex
}

func NewPersistent() *Persistent {
	return &Persistent{
		CurrentTerm: 0,
		VotedFor:    "",
		MachineLogs: []LogEntry{},
	}
}

func (p *Persistent) GetCopy() *Persistent {
	p.mu.RLock()
	defer p.mu.RUnlock()

	machineLogs := make([]LogEntry, len(p.MachineLogs))
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
	log.Println("before delete", p.MachineLogs, index, lastLogIndex, index <= lastLogIndex)
	if index <= lastLogIndex {
		p.MachineLogs = p.MachineLogs[:index]
		log.Println("after delete", p.MachineLogs, index)
	}
}

func (p *Persistent) AppendLogs(entries []LogEntry) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Copy existings logs
	logs := make([]LogEntry, len(p.MachineLogs), len(entries)+len(p.MachineLogs))
	copy(logs, p.MachineLogs)

	// Append new logs
	for _, entry := range entries {
		if entry.Term > 0 {
			logs = append(logs, entry)
		}
	}

	p.MachineLogs = logs
}

// Index starts at 1
func (p *Persistent) GetLogByIndex(index uint32) LogEntry {
	if index <= p.GetLastLogIndex() && index >= 1 {
		return p.MachineLogs[index-1]
	}
	return LogEntry{Term: 0}
}

// Returns index of last log
// Index starts at 1, return value 0 means there are no logs
func (p *Persistent) GetLastLogIndex() uint32 {
	return uint32(len(p.MachineLogs))
}

// Returns lastLogIndex and lastLog
func (p *Persistent) GetLastLog() (uint32, LogEntry) {
	lastLogIndex := p.GetLastLogIndex()
	lastLog := p.GetLogByIndex(lastLogIndex)
	return lastLogIndex, lastLog
}

func (p *Persistent) GetLastLogTerm() uint32 {
	_, lastLog := p.GetLastLog()
	return lastLog.Term
}

// Returns the entire logs stash from given index
func (p *Persistent) GetLogsFromIndex(index uint32) []LogEntry {
	logs := make([]LogEntry, 0, len(p.MachineLogs))
	lastLogIndex := p.GetLastLogIndex()
	for idx := index; idx <= lastLogIndex; idx++ {
		logs = append(logs, p.GetLogByIndex(idx))
	}
	return logs
}

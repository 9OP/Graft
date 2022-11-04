package entity

import (
	utils "graft/pkg/domain"
	"sync"

	log "github.com/sirupsen/logrus"
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
	if index <= p.GetLastLogIndex() {
		log.Info("DELETE LOGS FROM INDEX", index)
		// index-1 because index starts at 1
		p.MachineLogs = p.MachineLogs[:index-1]
	}
}

func (p *Persistent) AppendLogs(entries []LogEntry, prevLogIndex uint32) bool {
	// Should append only new entries
	if len(entries) == 0 {
		return false
	}

	/*
		Prevent appending logs twice
		Example:
		logs  -> [log1, log2, log3, log4, log5]
		term  ->   t1    t1    t2    t4    t4
		index ->   i1    i2    i3    i4    i5

		arguments -> [log3, log4, log5, log6], i2

		result -> [log1, log2, log3, log4, log5, log6]

		---
		We should increment prevLogIndex up to lastLogIndex.
		While incrementing a pointer in the entries slice.
		and then only copy the remaining "new" logs.

		prevLogIndex is necessay as it gives the offset of the
		entries argument, relative to p.MachineLogs
	*/

	p.mu.Lock()
	defer p.mu.Unlock()

	// Copy existings logs
	logs := make([]LogEntry, len(p.MachineLogs), len(entries)+len(p.MachineLogs))
	copy(logs, p.MachineLogs)

	// Find index of newLogs
	lastLogIndex := p.GetLastLogIndex()
	newLogsFromIndex := utils.Min(lastLogIndex-prevLogIndex, uint32(len(entries)))

	// Append new logs
	logs = append(logs, entries[newLogsFromIndex:]...)
	p.MachineLogs = logs

	log.Info("APPEND LOGS", len(entries[newLogsFromIndex:]))
	return len(entries[newLogsFromIndex:]) > 0
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

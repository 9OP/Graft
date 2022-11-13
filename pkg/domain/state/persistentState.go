package state

import (
	"errors"

	"graft/pkg/domain"
)

// https://github.com/golang/go/wiki/SliceTricks

type PersistentState struct {
	currentTerm uint32
	votedFor    string
	machineLogs []domain.LogEntry
}

func NewPersistentState(currentTerm uint32, votedFor string, machineLogs []domain.LogEntry) PersistentState {
	return PersistentState{
		currentTerm: currentTerm,
		votedFor:    votedFor,
		machineLogs: machineLogs,
	}
}

func NewDefaultPersistentState() PersistentState {
	return NewPersistentState(0, "", []domain.LogEntry{})
}

func (p PersistentState) CurrentTerm() uint32 {
	return p.currentTerm
}

func (p PersistentState) VotedFor() string {
	return p.votedFor
}

func (p PersistentState) LastLogIndex() uint32 {
	return uint32(len(p.machineLogs))
}

func (p PersistentState) LastLog() domain.LogEntry {
	lastLog, _ := p.MachineLog(p.LastLogIndex())
	return lastLog
}

func (p PersistentState) MachineLog(index uint32) (domain.LogEntry, error) {
	if index == 0 {
		return domain.LogEntry{}, nil
	}
	lastLogIndex := p.LastLogIndex()
	if index <= lastLogIndex {
		log := p.machineLogs[index-1]
		return log.Copy(), nil
	}
	return domain.LogEntry{}, errors.New("index out of range")
}

func (p PersistentState) MachineLogs() []domain.LogEntry {
	len := p.LastLogIndex()
	machineLogs := make([]domain.LogEntry, len)
	for i := uint32(1); i <= len; i += 1 {
		if log, err := p.MachineLog(i); err == nil {
			machineLogs[i-1] = log.Copy()
		}
	}
	return machineLogs
}

func (p PersistentState) MachineLogsFrom(index uint32) []domain.LogEntry {
	lastLogIndex := p.LastLogIndex()

	if index == 0 {
		logs := make([]domain.LogEntry, lastLogIndex)
		copy(logs, p.MachineLogs())
		return logs
	}

	if index <= lastLogIndex {
		logs := make([]domain.LogEntry, lastLogIndex-index+1)
		copy(logs, p.MachineLogs()[index-1:])
		return logs
	}

	return []domain.LogEntry{}
}

func (p PersistentState) WithCurrentTerm(term uint32) PersistentState {
	p.currentTerm = term
	return p
}

func (p PersistentState) WithVotedFor(vote string) PersistentState {
	p.votedFor = vote
	return p
}

func (p PersistentState) WithDeleteLogsFrom(index uint32) (PersistentState, bool) {
	// Delete logs from given index (include deletion)
	lastLogIndex := p.LastLogIndex()
	if index <= lastLogIndex && index >= 1 {
		// index-1 because index starts at 1
		logs := make([]domain.LogEntry, index-1)
		copy(logs, p.MachineLogs()[:index-1])
		p.machineLogs = logs
		return p, true
	}
	return p, false
}

func (p PersistentState) WithAppendLogs(prevLogIndex uint32, entries ...domain.LogEntry) (PersistentState, bool) {
	// Should append only new entries
	changed := false
	if len(entries) == 0 {
		return p, changed
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

	lastLogIndex := p.LastLogIndex()

	// Find index of newLogs
	newLogsFromIndex := lastLogIndex - prevLogIndex
	changed = len(entries[newLogsFromIndex:]) > 0

	if !changed {
		return p, changed
	}

	// Copy existings logs
	logs := make([]domain.LogEntry, lastLogIndex, uint32(len(entries))+lastLogIndex)

	// Append new logs
	logs = append(logs, entries[newLogsFromIndex:]...)
	copy(logs, p.MachineLogs())

	p.machineLogs = logs
	return p, changed
}

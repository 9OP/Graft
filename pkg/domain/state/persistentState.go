package state

import (
	"errors"

	"graft/pkg/domain"
	"graft/pkg/utils"
)

type PersistentState struct {
	currentTerm uint32
	votedFor    string
	// Important: index starts at 1
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
	lastLog, _ := p.Log(p.LastLogIndex())
	return lastLog
}

var errIndexOutOfRange = errors.New("index out of range")

func (p PersistentState) Log(index uint32) (domain.LogEntry, error) {
	if index == 0 {
		return domain.LogEntry{}, nil
	}
	if index <= p.LastLogIndex() {
		return p.machineLogs[index-1], nil
	}
	return domain.LogEntry{}, errIndexOutOfRange
}

// Return a slice copy of state machine logs
func (p PersistentState) Logs() []domain.LogEntry {
	machineLogs := make([]domain.LogEntry, len(p.machineLogs))
	copy(machineLogs, p.machineLogs)
	return machineLogs
}

func (p PersistentState) LogsFrom(index uint32) []domain.LogEntry {
	logs := p.Logs()
	if index == 0 {
		return logs
	}
	if index <= p.LastLogIndex() {
		return logs[index-1:]
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
	if index == 0 {
		p.machineLogs = []domain.LogEntry{}
		return p, true
	}
	logs := p.Logs()
	if index <= p.LastLogIndex() && index >= 1 {
		p.machineLogs = logs[:index-1]
		return p, true
	}
	p.machineLogs = logs
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
		entries argument, relative to p.Logs
	*/

	lastLogIndex := p.LastLogIndex()
	lenEntries := uint32(len(entries))

	// Find index of newLogs
	var newLogsFromIndex uint32
	if lastLogIndex >= prevLogIndex {
		newLogsFromIndex = utils.Min(lastLogIndex-prevLogIndex, lenEntries)
	} else {
		newLogsFromIndex = lenEntries
	}
	changed = len(entries[newLogsFromIndex:]) > 0

	if !changed {
		return p, changed
	}

	// Copy existings logs
	logs := make([]domain.LogEntry, lastLogIndex, lenEntries+lastLogIndex)

	// Append new logs
	logs = append(logs, entries[newLogsFromIndex:]...)
	copy(logs, p.Logs())

	p.machineLogs = logs
	return p, changed
}

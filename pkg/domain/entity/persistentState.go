package entity

import (
	"errors"
	utils "graft/pkg/domain"
)

// https://github.com/golang/go/wiki/SliceTricks

type PersistentState struct {
	currentTerm uint32
	votedFor    string
	machineLogs []LogEntry
}

func NewPersistentState(currentTerm uint32, votedFor string, machineLogs []LogEntry) PersistentState {
	return PersistentState{
		currentTerm: currentTerm,
		votedFor:    votedFor,
		machineLogs: machineLogs,
	}
}

func NewDefaultPersistentState() PersistentState {
	return NewPersistentState(0, "", []LogEntry{})
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

func (p PersistentState) LastLog() LogEntry {
	lastLog, _ := p.MachineLog(p.LastLogIndex())
	return lastLog
}

func (p PersistentState) MachineLog(index uint32) (LogEntry, error) {
	lastLogIndex := p.LastLogIndex()
	if index <= lastLogIndex && index >= 1 {
		log := p.machineLogs[index-1]
		return LogEntry{Term: log.Term, Value: log.Value}, nil
	}
	return LogEntry{}, errors.New("index out of range")
}

func (p PersistentState) MachineLogs() []LogEntry {
	len := p.LastLogIndex()
	machineLogs := make([]LogEntry, len)
	for i := uint32(1); i <= len; i += 1 {
		if log, err := p.MachineLog(i); err == nil {
			machineLogs[i-1] = log
		}
	}
	return machineLogs
}

func (p PersistentState) MachineLogsFrom(index uint32) []LogEntry {
	lastLogIndex := p.LastLogIndex()
	logs := make([]LogEntry, 0, lastLogIndex)
	if index <= lastLogIndex {
		copy(logs, p.machineLogs[index-1:])
	}
	return logs
}

func (p PersistentState) WithCurrentTerm(term uint32) PersistentState {
	p.currentTerm = term
	return p
}

func (p PersistentState) WithVotedFor(vote string) PersistentState {
	p.votedFor = vote
	return p
}

func (p PersistentState) WithDeleteLogsFrom(index uint32) PersistentState {
	// Delete logs from given index (include deletion)
	lastLogIndex := p.LastLogIndex()
	if index <= lastLogIndex && index >= 1 {
		logs := make([]LogEntry, 0, lastLogIndex)
		// index-1 because index starts at 1
		copy(logs, p.machineLogs[:index-1])
		p.machineLogs = logs
		return p
	}
	return p
}

func (p PersistentState) WithAppendLogs(prevLogIndex uint32, entries ...LogEntry) PersistentState {
	// Should append only new entries
	if len(entries) == 0 {
		return p
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

	numEntries := uint32(len(entries))
	lastLogIndex := p.LastLogIndex()

	// Copy existings logs
	logs := make([]LogEntry, lastLogIndex, numEntries+lastLogIndex)
	copy(logs, p.machineLogs)

	// Find index of newLogs
	newLogsFromIndex := utils.Min(lastLogIndex-prevLogIndex, numEntries)

	// Append new logs
	logs = append(logs, entries[newLogsFromIndex:]...)

	p.machineLogs = logs
	return p
}

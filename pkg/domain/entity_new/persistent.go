package entitynew

import (
	"errors"
	utils "graft/pkg/domain"
	"graft/pkg/domain/entity"
)

// https://github.com/golang/go/wiki/SliceTricks

type Persistent struct {
	currentTerm uint32
	votedFor    string
	machineLogs []entity.LogEntry
}

func NewPersistent() Persistent {
	return Persistent{
		currentTerm: 0,
		votedFor:    "",
		machineLogs: []entity.LogEntry{},
	}
}

func (p Persistent) CurrentTerm() uint32 {
	return p.currentTerm
}

func (p Persistent) VotedFor() string {
	return p.votedFor
}

func (p Persistent) LastLogIndex() uint32 {
	return uint32(len(p.machineLogs))
}

func (p Persistent) LastLog() entity.LogEntry {
	lastLog, _ := p.MachineLog(p.LastLogIndex())
	return lastLog
}

func (p Persistent) MachineLog(index uint32) (entity.LogEntry, error) {
	lastLogIndex := p.LastLogIndex()
	if index <= lastLogIndex && index >= 1 {
		log := p.machineLogs[index-1]
		return entity.LogEntry{Term: log.Term, Value: log.Value}, nil
	}
	return entity.LogEntry{}, errors.New("index out of range")
}

func (p Persistent) MachineLogsFrom(index uint32) []entity.LogEntry {
	lastLogIndex := p.LastLogIndex()
	logs := make([]entity.LogEntry, 0, lastLogIndex)
	if index <= lastLogIndex {
		copy(logs, p.machineLogs[index-1:])
	}
	return logs
}

func (p Persistent) WithCurrentTerm(term uint32) Persistent {
	p.currentTerm = term
	return p
}

func (p Persistent) WithVotedFor(vote string) Persistent {
	p.votedFor = vote
	return p
}

func (p Persistent) WithDeleteLogsFrom(index uint32) Persistent {
	// Delete logs from given index (include deletion)
	lastLogIndex := p.LastLogIndex()
	if index <= lastLogIndex && index >= 1 {
		logs := make([]entity.LogEntry, 0, lastLogIndex)
		// index-1 because index starts at 1
		copy(logs, p.machineLogs[:index-1])
		p.machineLogs = logs
		return p
	}
	return p
}

func (p Persistent) WithAppendLogs(entries []entity.LogEntry, prevLogIndex uint32) Persistent {
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
	logs := make([]entity.LogEntry, lastLogIndex, numEntries+lastLogIndex)
	copy(logs, p.machineLogs)

	// Find index of newLogs
	newLogsFromIndex := utils.Min(lastLogIndex-prevLogIndex, numEntries)

	// Append new logs
	logs = append(logs, entries[newLogsFromIndex:]...)

	p.machineLogs = logs
	return p
}

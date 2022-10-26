package entity

type MachineLog struct {
	Term  uint32 `json:"term"`
	Value string `json:"value"`
}

/*
 * Private entity available through state only
 */
type persistent struct {
	currentTerm uint32
	votedFor    string
	machineLogs []MachineLog
}

func NewPersistent() *persistent {
	return &persistent{
		currentTerm: 0,
		votedFor:    "",
		machineLogs: []MachineLog{{Term: 0, Value: "INIT"}},
	}
}

func (p *persistent) getPersistentCopy() *Persistent {
	machineLogs := make([]MachineLog, len(p.machineLogs))
	copy(machineLogs, p.machineLogs)

	return &Persistent{
		CurrentTerm: p.currentTerm,
		VotedFor:    p.votedFor,
		MachineLogs: machineLogs,
	}
}

func (p *persistent) SetVotedFor(id string) {
	p.votedFor = id
}
func (p *persistent) SetCurrentTerm(term uint32) {
	p.currentTerm = term
}
func (p *persistent) DeleteLogsFrom(index uint32) {
	lastLogIndex := int32(len(p.machineLogs))
	if index < uint32(lastLogIndex) {
		p.machineLogs = p.machineLogs[:index]
	}
}
func (p *persistent) AppendLogs(entries []string) {
	size := len(entries) + len(p.machineLogs)
	logs := make([]MachineLog, size)
	currentTerm := p.currentTerm

	for _, entry := range entries {
		newLog := MachineLog{Term: currentTerm, Value: entry}
		logs = append(logs, newLog)
	}
	p.machineLogs = logs
}

/*
 * Public entity
 */
type Persistent struct {
	CurrentTerm uint32       `json:"current_term"`
	VotedFor    string       `json:"voted_for"`
	MachineLogs []MachineLog `json:"machine_logs"`
}

func (p Persistent) GetLastLogIndex() uint32 {
	return uint32(len(p.MachineLogs))
}
func (p Persistent) GetLastLogTerm() uint32 {
	if lastLogIndex := p.GetLastLogIndex(); lastLogIndex != 0 {
		return p.MachineLogs[lastLogIndex-1].Term
	}
	return 0
}
func (p Persistent) GetLogByIndex(index uint32) MachineLog {
	if index < p.GetLastLogIndex() {
		// Return copy for safety
		return MachineLog{
			Term:  p.MachineLogs[index].Term,
			Value: p.MachineLogs[index].Value,
		}
	}
	return MachineLog{Term: 0}
}

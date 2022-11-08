package service

import (
	"fmt"
	"graft/pkg/domain/entity"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func (c ClusterNode) eval(entry string, entryType string) entity.EvalResult {
	cmd := exec.Command(
		c.fsmEval,
		entry,
		entryType,
		c.Id(),
		c.LeaderId(),
		c.Role().String(),
		fmt.Sprint(c.LastLogIndex()),
		fmt.Sprint(c.CurrentTerm()),
		fmt.Sprint(c.VotedFor()),
	)
	out, err := cmd.Output()
	// res :=
	log.Debugf("EVAL:\n\t%s\n\t%S\n\t%s", entry, string(out), err)
	return entity.EvalResult{
		Out: out,
		Err: err,
	}
}
func (c ClusterNode) init() {
	//
}

package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type fsm struct {
	path string
}

func NewFsm(path string) fsm {
	return fsm{path}
}

type fsmEvalType uint8

const (
	_ fsmEvalType = iota
	INIT
	QUERY
	COMMAND
)

type EvalInput struct {
	Id       string      `json:"id"`
	Data     string      `json:"data"`
	EvalType fsmEvalType `json:"type"`
}

func (f fsm) eval(input EvalInput) ExecuteOutput {
	var outb, errb bytes.Buffer
	payload, _ := json.Marshal(input)

	cmd := exec.Command(f.path, string(payload))
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	fmt.Println("out: ", outb.String())
	fmt.Println("err: ", errb.String(), err)
	if err != nil {
		return ExecuteOutput{
			Out: outb.Bytes(),
			Err: fmt.Errorf("fsm.eval error: %v", errb.String()),
		}
	}

	return ExecuteOutput{
		Out: outb.Bytes(),
		Err: nil,
	}
}

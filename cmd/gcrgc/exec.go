package main

import (
	"bytes"
	"os/exec"
)

// CmdExecutor can execute Cmd commands
type CmdExecutor interface {
	Exec(gcmd *Cmd) error
}

// Cmd represents a cli command
type Cmd struct {
	name   string
	args   []string
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

// NewCmd creates a new Cmd instance
func NewCmd(name string, args []string) *Cmd {
	var out, err bytes.Buffer

	return &Cmd{name, args, &out, &err}
}

// Cli represents a cli executor
type Cli struct {
}

// NewCli create a new Cli instance
func NewCli() *Cli {
	return &Cli{}
}

// Exec execute a gcloud command
// args are the arguments to be passed after the "gcloud" command
func (c Cli) Exec(gcmd *Cmd) error {
	cmd := exec.Command(gcmd.name, gcmd.args...)
	cmd.Stdout = gcmd.Stdout
	cmd.Stderr = gcmd.Stderr

	return cmd.Run()
}

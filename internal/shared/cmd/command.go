package cmd

import (
	"bytes"
	"context"
	"os/exec"
)

type Command struct {
	bin string
}

type Result struct {
	CommandLine string
	Stdout      string
	Stderr      string
}

func NewCommand(bin string) *Command {
	return &Command{
		bin: bin,
	}
}

func (c *Command) Run(ctx context.Context, args ...string) (*Result, error) {
	cmd := exec.CommandContext(ctx, c.bin, args...)

	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()

	return &Result{
		CommandLine: cmd.String(),
		Stdout:      outBuffer.String(),
		Stderr:      errBuffer.String(),
	}, err
}

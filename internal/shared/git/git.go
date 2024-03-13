package git

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Git struct {
	bin string
}

type cmdResult struct {
	stdout string
	stderr string
}

func NewGit(binary string) *Git {
	return &Git{
		bin: binary,
	}
}

type Commit struct {
	Message string
	Author  string
}

func (g *Git) Commit(ctx context.Context, commit *Commit) error {
	log.Printf("[git] commiting changes")

	args := []string{
		"commit",
		"-m",
		commit.Message,
	}

	if commit.Author != "" {
		args = append(args, fmt.Sprintf("--author=%s", commit.Author))
	}

	cmd := exec.CommandContext(ctx, g.bin, args...)
	if msg, err := cmd.Output(); err != nil {
		return fmt.Errorf("failed to execute %q: %w: %s", cmd.String(), err, msg)
	}

	log.Printf("[git] changes commited")

	return nil
}

func (g *Git) AddFile(ctx context.Context, filename string) error {
	log.Printf("[git] adding file %q", filename)

	cmd := exec.CommandContext(ctx, g.bin, "add", filename)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute %q: %w", cmd.String(), err)
	}

	log.Printf("[git] added file %q", filename)

	return nil
}

func (g *Git) Push(ctx context.Context) error {
	log.Printf("[git] pushing")

	cmd := exec.CommandContext(ctx, g.bin, "push")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute %q: %w", cmd.String(), err)
	}

	log.Printf("[git] pushed")

	return nil
}

func (g *Git) GetAddedAndModifiedFiles(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(
		ctx,
		g.bin,
		"diff",
		"--diff-filter=A",
		"--diff-filter=M",
		"--name-only",
		"HEAD",
	)
	res, err := g.run(cmd)
	if err != nil {
		return []string{}, fmt.Errorf("failed to execute %q: %w: %s", cmd.String(), err, res.stderr)
	}

	if res.stdout == "" {
		return []string{}, nil
	}

	return strings.Split(strings.Trim(res.stdout, "\n "), "\n"), nil
}

func (g *Git) run(cmd *exec.Cmd) (*cmdResult, error) {
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()

	return &cmdResult{
		stdout: outBuffer.String(),
		stderr: errBuffer.String(),
	}, err
}

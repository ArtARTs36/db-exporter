package git

import (
	"context"
	"fmt"
	"log"
	"os/exec"
)

type Git struct {
	bin string
}

func NewGit(binary string) *Git {
	return &Git{
		bin: binary,
	}
}

type Commit struct {
	Message   string
	FilePaths []string
}

func (g *Git) Commit(ctx context.Context, commit *Commit) error {
	if len(commit.FilePaths) > 0 {
		log.Printf("[git] adding files")

		for _, filename := range commit.FilePaths {
			err := g.AddFile(ctx, filename)
			if err != nil {
				return fmt.Errorf("failed to add file %q: %w", filename, err)
			}
		}

		log.Printf("[git] added %d files", len(commit.FilePaths))
	}

	return g.commit(ctx, commit)
}

func (g *Git) AddFile(ctx context.Context, filename string) error {
	cmd := exec.CommandContext(ctx, g.bin, "add", filename)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute %q: %w", cmd.String(), err)
	}

	log.Printf("[git] adding file %q", filename)

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

func (g *Git) commit(ctx context.Context, commit *Commit) error {
	log.Printf("[git] commiting changes")

	cmd := exec.CommandContext(ctx, g.bin, "commit", "-m", commit.Message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute %q: %w", cmd.String(), err)
	}

	log.Printf("[git] changes commited")

	return nil
}

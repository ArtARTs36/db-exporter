package git

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/cmd"
	"log/slog"
	"strings"
)

type Git struct {
	bin *cmd.Command
}

func NewGit(binary string) *Git {
	return &Git{
		bin: cmd.NewCommand(binary),
	}
}

type Commit struct {
	Message string
	Author  *Author
}

func (g *Git) Commit(ctx context.Context, commit *Commit) error {
	slog.InfoContext(ctx, "[git] committing changes")

	args := make([]string, 0)

	if commit.Author != nil {
		args = append(args, "-c")
		args = append(args, fmt.Sprintf("user.name=%s", commit.Author.Name))
		args = append(args, "-c")
		args = append(args, fmt.Sprintf("user.email=%s", commit.Author.Email))
	}

	args = append(args, "commit")
	args = append(args, "-m")
	args = append(args, commit.Message)

	if res, err := g.bin.Run(ctx, args...); err != nil {
		return fmt.Errorf("failed to execute %q: %w: %s", res.CommandLine, err, res.Stderr)
	}

	slog.InfoContext(ctx, "[git] changes committed")

	return nil
}

func (g *Git) AddFile(ctx context.Context, filename string) error {
	slog.InfoContext(ctx, fmt.Sprintf("[git] adding file %q", filename))

	if res, err := g.bin.Run(ctx, "add", filename); err != nil {
		return fmt.Errorf("failed to execute %q: %w: %s", res.CommandLine, err, res.Stderr)
	}

	slog.InfoContext(ctx, fmt.Sprintf("[git] added file %q", filename))

	return nil
}

func (g *Git) Push(ctx context.Context) error {
	slog.InfoContext(ctx, "[git] pushing")

	if res, err := g.bin.Run(ctx, "push"); err != nil {
		return fmt.Errorf("failed to execute %q: %w: %s", res.CommandLine, err, res.Stderr)
	}

	slog.InfoContext(ctx, "[git] pushed")

	return nil
}

func (g *Git) GetAddedAndModifiedFiles(ctx context.Context) ([]string, error) {
	args := []string{
		"diff",
		"--diff-filter=A",
		"--diff-filter=M",
		"--name-only",
		"HEAD",
	}
	res, err := g.bin.Run(ctx, args...)
	if err != nil {
		return []string{}, fmt.Errorf("failed to execute %q: %w: %s", res.CommandLine, err, res.Stderr)
	}

	if res.Stdout == "" {
		return []string{}, nil
	}

	return strings.Split(strings.Trim(res.Stdout, "\n "), "\n"), nil
}

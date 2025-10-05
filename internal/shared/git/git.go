package git

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/artarts36/db-exporter/internal/shared/cmd"
)

type Git struct {
	bin          *cmd.Command
	authorFinder AuthorFinder
}

func NewGit(binary string, authorFinder AuthorFinder) *Git {
	return &Git{
		bin:          cmd.NewCommand(binary),
		authorFinder: authorFinder,
	}
}

type Commit struct {
	Message string
	Author  *Author
}

func (g *Git) Commit(ctx context.Context, commit *Commit) error {
	author := commit.Author
	if author == nil {
		var err error
		author, err = g.authorFinder(ctx)
		if err != nil {
			return fmt.Errorf("failed to find author: %w", err)
		}
	}

	slog.
		With(slog.Any("commit", map[string]interface{}{
			"message": commit.Message,
			"author":  author,
		})).
		InfoContext(ctx, "[git] committing changes")

	args := make([]string, 0)

	if author != nil {
		args = append(args, "-c")
		args = append(args, fmt.Sprintf("user.name=%s", author.Name))
		args = append(args, "-c")
		args = append(args, fmt.Sprintf("user.email=%s", author.Email))
	}

	args = append(args, "commit")
	args = append(args, "-m")
	args = append(args, commit.Message)

	if res, err := g.bin.Run(ctx, args...); err != nil {
		return fmt.Errorf("failed to execute %q: %w: %s", res.CommandLine, err, res.Stderr)
	}

	slog.
		With(slog.Any("commit", commit)).
		InfoContext(ctx, "[git] changes committed")

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

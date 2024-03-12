package actions

import (
	"context"
	"fmt"

	"github.com/artarts36/db-exporter/internal/app/params"
	"github.com/artarts36/db-exporter/internal/shared/git"
)

type Commit struct {
	git *git.Git
}

func NewCommit(git *git.Git) *Commit {
	return &Commit{
		git: git,
	}
}

func (c *Commit) Supports(params *params.ActionParams) bool {
	return params.ExportParams.CommitMessage != "" || params.ExportParams.CommitPush
}

func (c *Commit) Run(ctx context.Context, params *params.ActionParams) error {
	err := c.git.Commit(ctx, &git.Commit{
		Message:   c.createCommitMessage(params),
		FilePaths: params.GeneratedFilesPaths,
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	if params.ExportParams.CommitPush {
		err = c.git.Push(ctx)
		if err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}
	}

	return nil
}

func (c *Commit) createCommitMessage(params *params.ActionParams) string {
	msg := params.ExportParams.CommitMessage

	if msg == "" {
		msg = "add documentation for database schema"
	}

	return msg
}

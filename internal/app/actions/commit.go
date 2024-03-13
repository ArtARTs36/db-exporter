package actions

import (
	"context"
	"fmt"
	"log/slog"
	"path"

	"strings"

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
	return params.ExportParams.CommitMessage != "" ||
		params.ExportParams.CommitAuthor != "" ||
		params.ExportParams.CommitPush
}

func (c *Commit) Run(ctx context.Context, params *params.ActionParams) error {
	err := c.checkUnexpectedFiles(ctx, params)
	if err != nil {
		return err
	}

	err = c.addFilesToGIt(ctx, params)
	if err != nil {
		return err
	}

	modifiedFiles, err := c.git.GetAddedAndModifiedFiles(ctx)
	if err != nil {
		return err
	}

	if len(modifiedFiles) == 0 {
		slog.DebugContext(ctx, "[commitaction] modified files not found, skip action")

		return nil
	}

	err = c.git.Commit(ctx, &git.Commit{
		Message: c.createCommitMessage(params),
		Author:  params.ExportParams.CommitAuthor,
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

func (c *Commit) addFilesToGIt(ctx context.Context, params *params.ActionParams) error {
	for _, f := range params.GeneratedFiles {
		err := c.git.AddFile(ctx, f.Path)
		if err != nil {
			return fmt.Errorf("failed to add file %q to git: %w", f, err)
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

func (c *Commit) checkUnexpectedFiles(ctx context.Context, params *params.ActionParams) error {
	modifiedFiles, err := c.git.GetAddedAndModifiedFiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to check modified files: %w", err)
	}

	if len(modifiedFiles) == 0 {
		return nil
	}

	modifiedFilesSet := map[string]bool{}
	for _, file := range modifiedFiles {
		modifiedFilesSet[path.Clean(file)] = true
	}

	for _, file := range params.GeneratedFiles {
		p := path.Clean(file.Path)

		if modifiedFilesSet[p] {
			delete(modifiedFilesSet, p)
		}
	}

	if len(modifiedFilesSet) == 0 {
		return nil
	}

	unexpectedModifiedFiles := make([]string, 0, len(modifiedFilesSet))
	for file := range modifiedFilesSet {
		unexpectedModifiedFiles = append(unexpectedModifiedFiles, file)
	}

	return fmt.Errorf(
		"git modified files contains unexpected files: %s",
		strings.Join(unexpectedModifiedFiles, ", "),
	)
}

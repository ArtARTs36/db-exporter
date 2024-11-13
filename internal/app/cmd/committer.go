package cmd

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"log/slog"
	"path"

	"strings"

	"github.com/artarts36/db-exporter/internal/shared/git"
)

type Committer struct {
	git *git.Git
}

func NewCommit(git *git.Git) *Committer {
	return &Committer{
		git: git,
	}
}

type commitParams struct {
	Commit         config.Commit
	GeneratedFiles []fs.FileInfo
}

func (c *Committer) Commit(ctx context.Context, params commitParams) error {
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
		slog.DebugContext(ctx, "[committer] modified files not found, skip action")

		return nil
	}

	var gitAuthor *git.Author
	if params.Commit.Author != "" {
		gitAuthor, err = git.NewAuthor(params.Commit.Author)
		if err != nil {
			return fmt.Errorf("failed to create git author: %w", err)
		}
	}

	err = c.git.Commit(ctx, &git.Commit{
		Message: c.createCommitMessage(params),
		Author:  gitAuthor,
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	if params.Commit.Push {
		err = c.git.Push(ctx)
		if err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}
	}

	return nil
}

func (c *Committer) addFilesToGIt(ctx context.Context, params commitParams) error {
	for _, f := range params.GeneratedFiles {
		slog.DebugContext(ctx, "[committer] adding file to git", slog.String("file", f.Path))

		err := c.git.AddFile(ctx, f.Path)
		if err != nil {
			return fmt.Errorf("failed to add file %q to git: %w", f, err)
		}
	}

	return nil
}

func (c *Committer) createCommitMessage(params commitParams) string {
	msg := params.Commit.Message

	if msg == "" {
		msg = "add documentation for database schema"
	}

	return msg
}

func (c *Committer) checkUnexpectedFiles(ctx context.Context, params commitParams) error {
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

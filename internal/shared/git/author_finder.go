package git

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type AuthorFinder func(ctx context.Context) (*Author, error)

var (
	githubActionsEnvs = []string{"GITHUB_ACTIONS", "DB_EXPORTER_USE_GITHUB_ACTIONS_BOT_COMMITTER"}

	githubActionsBotAuthor = &Author{
		Name:     "github-actions[bot]",
		Email:    "github-actions[bot]@users.noreply.github.com",
		Identity: "github-actions[bot] <github-actions[bot]@users.noreply.github.com>",
	}
)

func GithubActionsAuthorFinder() AuthorFinder {
	return returnAuthorWhenEnvLookup(githubActionsBotAuthor, githubActionsEnvs)
}

func returnAuthorWhenEnvLookup(author *Author, envs []string) AuthorFinder {
	return func(ctx context.Context) (*Author, error) {
		for _, env := range envs {
			if _, ok := os.LookupEnv(env); ok {
				return author, nil
			}
			slog.DebugContext(ctx, fmt.Sprintf("[git-author-finder] environment variable %s is missing", env))
		}
		return nil, nil
	}
}

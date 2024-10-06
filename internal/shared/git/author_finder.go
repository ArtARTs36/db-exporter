package git

import "os"

type AuthorFinder func() (*Author, error)

func GithubActionsAuthorFinder() AuthorFinder {
	return func() (*Author, error) {
		_, exists := os.LookupEnv("GITHUB_ACTIONS")
		if exists {
			return &Author{
				Name:     "github-actions[bot]",
				Email:    "github-actions[bot]@users.noreply.github.com",
				Identity: "github-actions[bot] <github-actions[bot]@users.noreply.github.com>",
			}, nil
		}
		return nil, nil
	}
}

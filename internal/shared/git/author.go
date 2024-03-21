package git

import (
	"fmt"
	"regexp"
	"strings"
)

type Author struct {
	Name     string
	Email    string
	Identity string
}

func NewAuthor(identity string) (*Author, error) {
	const identityParts = 3

	var re = regexp.MustCompile(`(?m)(.*)<(.*)>`)

	matches := re.FindAllStringSubmatch(identity, identityParts)
	if len(matches) == 0 || len(matches[0]) < identityParts {
		return nil, fmt.Errorf("failed to parse author identity: %s", identity)
	}

	name := strings.Trim(matches[0][1], " ")
	if name == "" {
		return nil, fmt.Errorf("failed to parse author identity: name is empty")
	}

	return &Author{
		Name:     name,
		Email:    matches[0][2],
		Identity: identity,
	}, nil
}

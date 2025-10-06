package dbml

import (
	"fmt"
	"strings"
)

type RelationSubject struct {
	Table  string
	Column string
}

func ParseRelationSubject(subject string) (*RelationSubject, error) {
	parts := strings.Split(subject, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("failed to parse relation subject: given %d dot parts, expected: 2", len(parts))
	}

	return &RelationSubject{
		Table:  parts[0],
		Column: parts[1],
	}, nil
}

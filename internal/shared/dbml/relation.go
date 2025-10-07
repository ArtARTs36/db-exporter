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
	const partsCount = 2

	parts := strings.Split(subject, ".")
	if len(parts) < partsCount {
		return nil, fmt.Errorf("failed to parse relation subject: given %d dot parts, expected: 2", len(parts))
	}

	return &RelationSubject{
		Table:  parts[0],
		Column: parts[1],
	}, nil
}

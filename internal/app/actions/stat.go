package actions

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/artarts36/db-exporter/internal/app/params"
	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type Stat struct {
	fs fs.Driver
}

func NewStat(fs fs.Driver) *Stat {
	return &Stat{
		fs: fs,
	}
}

func (*Stat) Supports(params *params.ActionParams) bool {
	return params.ExportParams.Stat
}

func (c *Stat) Run(_ context.Context, params *params.ActionParams) error {
	maxPathLen := 0

	for _, path := range params.GeneratedFilesPaths {
		if len(path) > maxPathLen {
			maxPathLen = len(path)
		}
	}

	fmt.Printf("file%s| state\n", strings.Repeat(" ", maxPathLen-3))
	fmt.Println(strings.Repeat("-", maxPathLen+10))

	for _, path := range params.GeneratedFilesPaths {
		state := "created"

		fileStat, err := c.fs.Stat(path)
		if err != nil {
			log.Printf("[stataction] failed to stat %q: %s", path, err)
			state = "unknown"
		} else if params.StartedAt.After(fileStat.CreatedAt) {
			state = "updated"
		}

		fmt.Printf("%s%s| %s\n", path, strings.Repeat(" ", maxPathLen-len(path)+1), state)
	}

	return nil
}

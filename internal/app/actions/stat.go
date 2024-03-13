package actions

import (
	"context"
	"github.com/artarts36/db-exporter/internal/app/params"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"log"
)

type Stat struct {
	fs         fs.Driver
	tblPrinter tablePrinter
}

type tablePrinter func(headers []string, rows [][]string)

func NewStat(fs fs.Driver, tblPrinter tablePrinter) *Stat {
	return &Stat{
		fs:         fs,
		tblPrinter: tblPrinter,
	}
}

func (*Stat) Supports(params *params.ActionParams) bool {
	return params.ExportParams.Stat
}

func (c *Stat) Run(_ context.Context, params *params.ActionParams) error {
	rows := make([][]string, 0, len(params.GeneratedFilesPaths))

	for _, path := range params.GeneratedFilesPaths {
		state := "created"

		fileStat, err := c.fs.Stat(path)
		if err != nil {
			log.Printf("[stataction] failed to stat %q: %s", path, err)
			state = "unknown"
		} else if params.StartedAt.After(fileStat.CreatedAt) {
			state = "updated"
		}

		rows = append(rows, []string{
			path,
			state,
		})
	}

	c.tblPrinter([]string{"file", "state"}, rows)

	return nil
}

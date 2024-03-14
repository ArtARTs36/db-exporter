package actions

import (
	"context"
	"fmt"

	"github.com/artarts36/db-exporter/internal/app/params"
)

type Stat struct {
	tblPrinter tablePrinter
}

type tablePrinter func(headers []string, rows [][]string)

func NewStat(tblPrinter tablePrinter) *Stat {
	return &Stat{
		tblPrinter: tblPrinter,
	}
}

func (*Stat) Supports(params *params.ActionParams) bool {
	return params.ExportParams.Stat
}

func (c *Stat) Run(_ context.Context, params *params.ActionParams) error {
	rows := make([][]string, 0, len(params.GeneratedFiles))

	for _, file := range params.GeneratedFiles {
		rows = append(rows, []string{
			file.Path,
			fmt.Sprintf("%d", file.Size),
		})
	}

	c.tblPrinter([]string{"file", "size"}, rows)

	return nil
}

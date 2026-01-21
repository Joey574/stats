package stats

import (
	"fmt"
	"strings"

	"github.com/Joey574/stats/internal/table"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

type CompiledTable struct {
	*table.Table
}

type Record struct {
	*table.Record
	Stddev float64
	Mean   float64
	Sem    float64
	CI95   float64
	CV     float64
}

func (c *CompiledTable) Compile() {
	r := ColumnStats(c.Table)
	RowStats(c.Table)
	c.Rows = append(c.Rows, r...)
}

func (c *CompiledTable) Dump(renderer tw.Renderer) string {
	var b strings.Builder
	rows, cols := c.Size()

	b.WriteString(fmt.Sprintf("Table: \"%s\" (%d x %d)\n", c.Name, rows, cols))

	writer := tablewriter.NewTable(&b,
		tablewriter.WithRenderer(renderer))
	writer.Header(c.Keys)

	for _, r := range c.Rows {
		writer.Append(r.Compose())
	}

	writer.Render()
	return b.String()
}

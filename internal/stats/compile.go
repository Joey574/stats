package stats

import (
	"fmt"
	"strings"

	"github.com/Joey574/stats/internal/table"
	"github.com/olekukonko/tablewriter"
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
	RowStats(c.Table)
	c.Rows = append(c.Rows, ColumnStats(c.Table)...)
}

func (c *CompiledTable) Dump() string {
	var b strings.Builder
	rows, cols := c.Size()

	b.WriteString(fmt.Sprintf("Table: \"%s\" (%d x %d)\n", c.Name, rows, cols))

	colTable := tablewriter.NewWriter(&b)
	colTable.Header(c.Headers())

	for _, r := range c.Rows {
		colTable.Append(r.Row())
	}

	colTable.Render()
	return b.String()
}

package stats

import (
	"github.com/Joey574/stats/internal/table"
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

func (c *CompiledTable) CompileDataTable() {
	r := ColumnStats(c.Table)
	RowStats(c.Table)
	c.Rows = append(c.Rows, r...)
}

func (c *CompiledTable) CompileForceTable() {
}

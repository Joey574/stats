package stats

import (
	"github.com/Joey574/stats/internal/table"
)

type CompiledTable struct {
	*table.Table
}

func (c *CompiledTable) CompileDataTable() {
	TableStats(c.Table)
}

func (c *CompiledTable) CompileForceTable() {
}

package stats

import (
	"github.com/Joey574/stats/internal/table"
)

type CompiledTable struct {
	*table.Table
}

func (c *CompiledTable) CompileDataTable() {
	cols := len(c.Keys)

	if cols > 1 {
		rowStats := tableStats(c.Table)
		c.Keys = append(c.Keys, []string{"MEAN", "STDDEV", "SEM", "CI₉₅", "CV"}...)
		for i := range c.Rows {
			c.Rows[i].Values = append(c.Rows[i].Values,
				table.Value{X: rowStats[i*nMetrics+0]},
				table.Value{X: rowStats[i*nMetrics+1]},
				table.Value{X: rowStats[i*nMetrics+2]},
				table.Value{X: rowStats[i*nMetrics+3], Prefix: "±"},
				table.Value{X: rowStats[i*nMetrics+4], Suffix: "%"},
			)
		}
	}
}

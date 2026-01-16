package main

import (
	"fmt"
	"strings"

	"github.com/Joey574/stats/internal/stats"
	"github.com/Joey574/stats/internal/table"
	"github.com/olekukonko/tablewriter"
)

type CompiledCol struct {
	Name string
	stats.Data
}

type CompiledRow struct {
	table.Row
	stats.RowData
}

type CompiledTable struct {
	Table *table.Table
	Cols  []*CompiledCol
	Rows  []*CompiledRow
}

func (t *CompiledTable) HasRowStatData() bool {
	for _, r := range t.Rows {
		if r.Diff != nil || r.Error != nil {
			return true
		}
	}
	return false
}

func (t *CompiledTable) Dump() string {
	var b strings.Builder

	// --- Metadata/Column Stats Table ---
	b.WriteString(fmt.Sprintf("Table: %s\n", t.Table.Name))
	colTable := tablewriter.NewWriter(&b)
	colTable.Header([]string{"Key", "Mean", "Stddev", "SEM", "CI95", "CV"})

	for _, k := range t.Cols {
		colTable.Append([]string{
			k.Name,
			fmt.Sprintf("%.2f", k.Mean),
			fmt.Sprintf("%.2f", k.Stddev),
			fmt.Sprintf("%.2f", k.Sem),
			fmt.Sprintf("±%.2f", k.CI95),
			fmt.Sprintf("%.2f%%", k.CV),
		})
	}
	colTable.Render()

	if t.HasRowStatData() {
		// --- Per Row Table ---
		b.WriteString(fmt.Sprintf("\nTable: %s\n", t.Table.Name))
		rowTable := tablewriter.NewWriter(&b)
		rowTable.Header([]string{"Idx", "Diff", "Error"})

		for i, r := range t.Rows {
			diffStr, errStr := "N/A", "N/A"
			if r.Diff != nil {
				diffStr = fmt.Sprintf("%.2f%%", *r.Diff)
			}
			if r.Error != nil {
				errStr = fmt.Sprintf("%.2f%%", *r.Error)
			}
			rowTable.Append([]string{fmt.Sprintf("%d", i), diffStr, errStr})
		}
		rowTable.Render()

	}

	return b.String()
}

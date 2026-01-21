package stats

import (
	"testing"

	"github.com/Joey574/stats/internal/table"
)

func buildTestTable() CompiledTable {
	rows := []*table.Record{
		{
			Label: "joey",
			Units: "s",
			Values: []table.Value{
				{X: 2.22},
				{X: 2.74},
				{X: 3.31},
				{X: 3.05},
				{X: 3.11},
				{X: 2.99},
				{X: 2.99},
				{X: 2.61},
				{X: 2.99},
				{X: 3.42},
			},
		},
		{
			Label: "sujan",
			Units: "s",
			Values: []table.Value{
				{X: 3.12},
				{X: 2.89},
				{X: 2.91},
				{X: 2.99},
				{X: 2.91},
				{X: 2.72},
				{X: 2.92},
				{X: 3.10},
				{X: 2.89},
				{X: 2.99},
			},
		},
		{
			Label: "phillip",
			Units: "s",
			Values: []table.Value{
				{X: 2.76},
				{X: 3.10},
				{X: 2.54},
				{X: 2.57},
				{X: 2.57},
				{X: 2.89},
				{X: 2.52},
				{X: 2.49},
				{X: 2.71},
				{X: 3.11},
			},
		},
		{
			Label: "lamine",
			Units: "s",
			Values: []table.Value{
				{X: 1.99},
				{X: 2.41},
				{X: 2.95},
				{X: 2.52},
				{X: 2.31},
				{X: 2.02},
				{X: 2.39},
				{X: 2.97},
				{X: 2.34},
				{X: 2.27},
			},
		},
	}

	return CompiledTable{
		&table.Table{
			Name: "Test",
			Keys: []string{"label", "x1", "x2", "x3", "x4", "x5", "x6", "x7", "x8", "x9", "x10"},
			Rows: rows,
		},
	}

}

func BenchmarkCompileDataTable(b *testing.B) {
	b.ReportAllocs()

	tables := make([]CompiledTable, b.N)
	for i := 0; i < b.N; i++ {
		tables[i] = buildTestTable()
	}
	b.SetBytes(tables[0].Bytes())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tables[i].CompileDataTable()
	}
}

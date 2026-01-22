package stats

import (
	"testing"

	"github.com/Joey574/stats/internal/table"
)

func BenchmarkTableStats(b *testing.B) {
	b.ReportAllocs()

	tables := make([]CompiledTable, b.N)
	for i := 0; i < b.N; i++ {
		t, err := table.ParseTestTable("testdata/csv/test1.csv")
		if err != nil {
			b.Fatal(err)
		}

		tables[i] = CompiledTable{
			Table: &t,
		}
	}
	b.SetBytes(tables[0].Bytes())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tableStats(tables[i].Table)
	}
}

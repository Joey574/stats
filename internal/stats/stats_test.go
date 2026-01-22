package stats

import "testing"

func BenchmarkTableStats(b *testing.B) {
	b.ReportAllocs()

	tables := make([]CompiledTable, b.N)
	for i := 0; i < b.N; i++ {
		tables[i] = buildTestTable()
	}
	b.SetBytes(tables[0].Bytes())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tableStats(tables[i].Table)
	}
}

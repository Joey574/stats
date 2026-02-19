package stats

import (
	"testing"

	"github.com/Joey574/stats/internal/parser"
	fixtures "github.com/Joey574/stats/internal/testfixtures"
)

func BenchmarkCompileDataTable(bench *testing.B) {
	test := func(b *testing.B, path string) {
		b.StopTimer()
		t, err := parser.ParseTestTable(path)
		if err != nil {
			b.Fatal(err)
		}
		ct := CompiledTable{Table: t}

		b.StartTimer()
		ct.CompileDataTable()
	}

	fixtures.TestAgainstCSV(bench, test)
}

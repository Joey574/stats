package stats

import (
	"testing"

	"github.com/Joey574/stats/internal/parser"
	fixtures "github.com/Joey574/stats/internal/testfixtures"
)

func BenchmarkTableStats(bench *testing.B) {
	test := func(b *testing.B, path string) {
		b.StopTimer()
		t, err := parser.ParseTestTable(path)
		if err != nil {
			b.Fatal(err)
		}

		b.StartTimer()
		tableStats(t)
	}

	fixtures.TestAgainstCSV(bench, test)
}

package stats

import (
	"testing"

	"github.com/Joey574/stats/internal/table"
	fixtures "github.com/Joey574/stats/internal/testfixtures"
)

func BenchmarkTableStats(bench *testing.B) {
	test := func(b *testing.B, path string) {
		b.StopTimer()
		t, err := table.ParseTestTable(path)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
		tableStats(&t)
	}

	fixtures.TestAgainstCSV(bench, test)
}

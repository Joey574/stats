package table

import (
	"testing"

	fixtures "github.com/Joey574/stats/internal/testfixtures"
)

func BenchmarkParseTable(bench *testing.B) {
	test := func(b *testing.B, path string) {
		_, err := ParseTestTable(path)
		if err != nil {
			b.Fatal(err)
		}
	}

	fixtures.TestAgainstCSV(bench, test)
}

package table

import (
	"testing"

	fixtures "github.com/Joey574/stats/internal/testfixtures"
)

func BenchmarkParseTestTable(b *testing.B) {
	var sink Table
	var err error
	b.ReportAllocs()
	b.ResetTimer()

	entries, err := fixtures.TestCSV.ReadDir("testdata/csv")
	if err != nil {
		b.Fatal(err)
	}

	for _, entry := range entries {
		b.Run(entry.Name(), func(b *testing.B) {
			for b.Loop() {
				sink, err = ParseTestTable("testdata/csv/" + entry.Name())
				if err != nil {
					b.Fatal(err)
				}
			}

			b.SetBytes(sink.Bytes())
		})
	}
}

package fixtures

import (
	"testing"
)

const dir = "testdata/csv"

func TestAgainstCSV(bench *testing.B, test func(b *testing.B, path string)) {
	entries, err := TestCSV.ReadDir(dir)
	if err != nil {
		bench.Fatal(err)
	}

	for _, entry := range entries {
		bench.Run(entry.Name(), func(b *testing.B) {
			info, _ := entry.Info()
			b.SetBytes(info.Size())
			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				test(b, dir+"/"+entry.Name())
			}
		})
	}
}

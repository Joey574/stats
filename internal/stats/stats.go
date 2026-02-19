package stats

import (
	"math"

	"github.com/Joey574/stats/internal/table"
)

const nMetrics = 5
const nilValue = math.SmallestNonzeroFloat64

func TableStats(t *table.Table) {
	cols := len(t.Keys)
	rowStats := tableStats(t)

	// append row stats if n cols > 1
	if cols > 1 {
		t.Keys = append(t.Keys, []string{"MEAN", "STDDEV", "SEM", "CI₉₅", "CV"}...)
		for i := range t.Rows {
			t.Rows[i].Values = append(t.Rows[i].Values,
				table.Value{X: rowStats[0][i]},
				table.Value{X: rowStats[1][i]},
				table.Value{X: rowStats[2][i]},
				table.Value{X: rowStats[3][i], Prefix: "±"},
				table.Value{X: rowStats[4][i], Suffix: "%"},
			)
		}
	}
}

func tableStats(t *table.Table) [nMetrics][]float64 {
	rowCount := len(t.Rows)
	rowBlock := make([]float64, nMetrics*rowCount)
	rowFilled := make([]uint32, rowCount)

	// parse out row blocked data
	rowMeans, rowStds, rowSems, rowCi95, rowCvs :=
		rowBlock[:rowCount],
		rowBlock[rowCount:2*rowCount],
		rowBlock[2*rowCount:3*rowCount],
		rowBlock[3*rowCount:4*rowCount],
		rowBlock[4*rowCount:5*rowCount]

	meanTableStats(t, rowMeans, rowFilled)
	stddevTableStats(t, rowStds, rowMeans, rowFilled)
	miscTableStats(rowMeans, rowStds, rowSems, rowCi95, rowCvs, rowFilled)

	return [nMetrics][]float64{rowMeans, rowStds, rowSems, rowCi95, rowCvs}
}

func meanTableStats(t *table.Table, rm []float64, rf []uint32) {
	// collect sums and filled count
	for i, r := range t.Rows {
		for k := range r.Values {
			if r.Values[k].X == nilValue {
				continue
			}

			rm[i] += r.Values[k].X
			rf[i]++
		}
	}

	for i := range rm {
		rm[i] /= float64(rf[i])
	}
}

func stddevTableStats(t *table.Table, rs []float64, rm []float64, rf []uint32) {
	// compute variance data
	for i, r := range t.Rows {
		for k := range r.Values {
			if r.Values[k].X == nilValue {
				continue
			}

			rv := r.Values[k].X - rm[i]
			rs[i] += rv * rv
		}
	}

	// compute row stds
	for i := range rs {
		rs[i] = math.Sqrt(rs[i] / float64(rf[i]))
	}
}

func miscTableStats(m []float64, std []float64, se []float64, ci []float64, cv []float64, f []uint32) {
	// loop to compute col SEMs, CVs, and CI95s
	for i := range m {
		se[i] = std[i] / math.Sqrt(float64(f[i]))
		cv[i] = std[i] / m[i] * 100.0
		ci[i] = 1.96 * se[i]
	}
}

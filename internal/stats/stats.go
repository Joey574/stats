package stats

import (
	"math"
	"slices"

	"github.com/Joey574/stats/internal/table"
)

const nilString = ""
const nilValue = math.SmallestNonzeroFloat64

func TableStats(t *table.Table) {
	units := t.Rows[0].Units
	cols, rows := tableStats(t)

	conv := func(x []float64, prefix string, suffix string, usesUnits bool) []table.Value {
		val := make([]table.Value, len(x))

		for i := range x {
			val[i] = table.Value{
				X:         x[i],
				UsesUnits: usesUnits,
				Prefix:    prefix,
				Suffix:    suffix,
			}
		}
		return val
	}

	// append row stats
	t.Keys = append(t.Keys, []string{"MEAN", "STDDEV", "SEM", "CI₉₅", "CV"}...)
	for i, r := range t.Rows {
		r.Values = slices.Grow(r.Values, 5)
		r.Values = append(r.Values,
			table.Value{X: rows[0][i], UsesUnits: true},
			table.Value{X: rows[1][i], UsesUnits: true},
			table.Value{X: rows[2][i], UsesUnits: true},
			table.Value{X: rows[3][i], UsesUnits: true, Prefix: "±"},
			table.Value{X: rows[4][i], Suffix: "%"},
		)
	}

	// append column stats
	t.Rows = slices.Grow(t.Rows, 5)
	t.Rows = append(t.Rows, []*table.Record{
		{
			Label:  "MEAN",
			Units:  units,
			Values: conv(cols[0], nilString, nilString, true),
		},
		{
			Label:  "STDDEV",
			Units:  units,
			Values: conv(cols[1], nilString, nilString, true),
		},
		{
			Label:  "SEM",
			Units:  units,
			Values: conv(cols[2], nilString, nilString, true),
		},
		{
			Label:  "CI₉₅",
			Units:  units,
			Values: conv(cols[3], "±", nilString, true),
		},
		{
			Label:  "CV",
			Units:  nilString,
			Values: conv(cols[4], nilString, "%", false),
		},
	}...)
}

func tableStats(t *table.Table) ([5][]float64, [5][]float64) {
	rowCount := len(t.Rows)
	colCount := len(t.Rows[0].Values)

	rowBlock := make([]float64, 5*rowCount)
	colBlock := make([]float64, 5*colCount)

	rowMeans, rowStds, rowSems, rowCi95, rowCvs := rowBlock[:rowCount], rowBlock[rowCount:2*rowCount], rowBlock[2*rowCount:3*rowCount], rowBlock[3*rowCount:4*rowCount], rowBlock[4*rowCount:5*rowCount]
	colMeans, colStds, colSems, colCi95, colCvs := colBlock[:colCount], colBlock[colCount:2*colCount], colBlock[2*colCount:3*colCount], colBlock[3*colCount:4*colCount], colBlock[4*colCount:5*colCount]

	colFilled := make([]uint32, colCount)
	rowFilled := make([]uint32, rowCount)

	// collect sums into mean
	for i, r := range t.Rows {
		for k := range r.Values {
			if r.Values[k].X == nilValue {
				continue
			}

			colMeans[k] += r.Values[k].X
			rowMeans[i] += r.Values[k].X
			colFilled[k]++
			rowFilled[i]++
		}
	}

	// compute col mean
	for i := range colMeans {
		colMeans[i] /= float64(colFilled[i])
	}

	for i := range rowMeans {
		rowMeans[i] /= float64(rowFilled[i])
	}

	// compute variance
	for i, r := range t.Rows {
		for k := range r.Values {
			if r.Values[k].X == nilValue {
				continue
			}

			cs := r.Values[k].X - colMeans[k]
			rs := r.Values[k].X - rowMeans[i]

			colStds[k] += cs * cs
			rowStds[i] += rs * rs
		}
	}

	// compute col stds
	for i := range colStds {
		colStds[i] = math.Sqrt(colStds[i] / float64(colFilled[i]))
	}

	// compute row stds
	for i := range rowStds {
		rowStds[i] = math.Sqrt(rowStds[i] / float64(rowFilled[i]))
	}

	// loop to compute col SEMs, CVs, and CI95s
	for i := range colSems {
		colSems[i] = colStds[i] / math.Sqrt(float64(colFilled[i]))
		colCvs[i] = colStds[i] / colMeans[i] * 100.0
		colCi95[i] = 1.96 * colSems[i]
	}

	// loop to compute row SEMs, CVs, and CI95s
	for i := range rowSems {
		rowSems[i] = rowStds[i] / math.Sqrt(float64(rowFilled[i]))
		rowCvs[i] = rowStds[i] / rowMeans[i] * 100.0
		rowCi95[i] = 1.96 * rowSems[i]
	}

	return [5][]float64{colMeans, colStds, colSems, colCi95, colCvs}, [5][]float64{rowMeans, rowStds, rowSems, rowCi95, rowCvs}
}

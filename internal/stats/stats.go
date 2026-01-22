package stats

import (
	"math"

	"github.com/Joey574/stats/internal/table"
)

const nMetrics = 5
const nilString = ""
const nilValue = math.SmallestNonzeroFloat64

func TableStats(t *table.Table) {
	cols, rows := len(t.Keys), len(t.Rows)

	units := t.Rows[0].Units
	colStats, rowStats := tableStats(t)

	valuePool := make([]table.Value, nMetrics*cols)
	conv := func(idx int, prefix string, suffix string, usesUnits bool) []table.Value {
		c := colStats[idx]
		for i := range c {
			valuePool[idx*cols+i] = table.Value{
				X:         c[i],
				UsesUnits: usesUnits,
				Prefix:    prefix,
				Suffix:    suffix,
			}
		}
		return valuePool[idx*cols : idx*cols+cols]
	}

	// append row stats if n cols > 1
	if cols > 1 {
		t.Keys = append(t.Keys, []string{"MEAN", "STDDEV", "SEM", "CI₉₅", "CV"}...)
		for i, r := range t.Rows {
			r.Values = append(r.Values,
				table.Value{X: rowStats[0][i], UsesUnits: true},
				table.Value{X: rowStats[1][i], UsesUnits: true},
				table.Value{X: rowStats[2][i], UsesUnits: true},
				table.Value{X: rowStats[3][i], UsesUnits: true, Prefix: "±"},
				table.Value{X: rowStats[4][i], Suffix: "%"},
			)
		}
	}

	// append column stats if n rows > 1
	if rows > 1 {
		t.Rows = append(t.Rows, []*table.Record{
			{
				Label:  "MEAN",
				Units:  units,
				Values: conv(0, nilString, nilString, true),
			},
			{
				Label:  "STDDEV",
				Units:  units,
				Values: conv(1, nilString, nilString, true),
			},
			{
				Label:  "SEM",
				Units:  units,
				Values: conv(2, nilString, nilString, true),
			},
			{
				Label:  "CI₉₅",
				Units:  units,
				Values: conv(3, "±", nilString, true),
			},
			{
				Label:  "CV",
				Units:  nilString,
				Values: conv(4, nilString, "%", false),
			},
		}...)
	}
}

func tableStats(t *table.Table) ([nMetrics][]float64, [nMetrics][]float64) {
	rowCount := len(t.Rows)
	colCount := len(t.Keys)

	dataBlock := make([]float64, nMetrics*(rowCount+colCount))
	rowBlock := dataBlock[:nMetrics*rowCount]
	colBlock := dataBlock[nMetrics*rowCount:]

	// parse out row blocked data
	rowMeans, rowStds, rowSems, rowCi95, rowCvs :=
		rowBlock[:rowCount],
		rowBlock[rowCount:2*rowCount],
		rowBlock[2*rowCount:3*rowCount],
		rowBlock[3*rowCount:4*rowCount],
		rowBlock[4*rowCount:5*rowCount]

	// parse out col blocked data
	colMeans, colStds, colSems, colCi95, colCvs :=
		colBlock[:colCount],
		colBlock[colCount:2*colCount],
		colBlock[2*colCount:3*colCount],
		colBlock[3*colCount:4*colCount],
		colBlock[4*colCount:5*colCount]

	// parse out more blocked data
	colFilled := make([]uint32, colCount+rowCount)
	rowFilled := colFilled[colCount:]
	colFilled = colFilled[:colCount]

	meanTableStats(t, colMeans, rowMeans, colFilled, rowFilled)
	stddevTableStats(t, colStds, rowStds, colMeans, rowMeans, colFilled, rowFilled)

	miscTableStats(colMeans, colStds, colSems, colCi95, colCvs, colFilled)
	miscTableStats(rowMeans, rowStds, rowSems, rowCi95, rowCvs, rowFilled)

	return [nMetrics][]float64{colMeans, colStds, colSems, colCi95, colCvs}, [nMetrics][]float64{rowMeans, rowStds, rowSems, rowCi95, rowCvs}
}

func meanTableStats(t *table.Table, cm []float64, rm []float64, cf []uint32, rf []uint32) {
	// collect sums into mean
	for i, r := range t.Rows {
		for k := range r.Values {
			if r.Values[k].X == nilValue {
				continue
			}

			cm[k] += r.Values[k].X
			rm[i] += r.Values[k].X
			cf[k]++
			rf[i]++
		}
	}

	for i := range cm {
		cm[i] /= float64(cf[i])
	}

	for i := range rm {
		rm[i] /= float64(rf[i])
	}
}

func stddevTableStats(t *table.Table, cs []float64, rs []float64, cm []float64, rm []float64, cf []uint32, rf []uint32) {
	for i, r := range t.Rows {
		for k := range r.Values {
			if r.Values[k].X == nilValue {
				continue
			}

			cv := r.Values[k].X - cm[k]
			rv := r.Values[k].X - rm[i]

			cs[k] += cv * cv
			rs[i] += rv * rv
		}
	}

	// compute col stds
	for i := range cs {
		cs[i] = math.Sqrt(cs[i] / float64(cf[i]))
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

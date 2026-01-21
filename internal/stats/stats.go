package stats

import (
	"math"
	"sync"

	"github.com/Joey574/stats/internal/table"
)

const nilString = ""
const nilValue = math.SmallestNonzeroFloat64

func ColumnStats(t *table.Table) []*table.Record {
	units := t.Rows[0].Units

	avgs := columnAverages(t)
	stds := columnStddevs(t, avgs)
	sems := columnSEMs(t, stds)
	ci95s := columnCI95s(sems)
	cvs := columnCVs(stds, avgs)

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

	return []*table.Record{
		{
			Label:  "MEAN",
			Units:  units,
			Values: conv(avgs, nilString, nilString, true),
		},
		{
			Label:  "STDDEV",
			Units:  units,
			Values: conv(stds, nilString, nilString, true),
		},
		{
			Label:  "SEM",
			Units:  units,
			Values: conv(sems, nilString, nilString, true),
		},
		{
			Label:  "CI95",
			Units:  units,
			Values: conv(ci95s, "±", nilString, true),
		},
		{
			Label:  "CV",
			Units:  nilString,
			Values: conv(cvs, nilString, "%", false),
		},
	}
}

func columnAverages(t *table.Table) []float64 {
	// gather sums into avgs
	avgs := make([]float64, len(t.Rows[0].Values))
	filled := make([]uint32, len(avgs))

	for _, r := range t.Rows {
		for i := range r.Values {
			if r.Values[i].X == nilValue {
				continue
			}

			avgs[i] += r.Values[i].X
			filled[i]++
		}
	}

	// average out based on supplied values
	for i := range avgs {
		avgs[i] /= float64(filled[i])
	}

	return avgs
}

func columnStddevs(t *table.Table, mean []float64) []float64 {
	// gather sums into variance
	variance := make([]float64, len(mean))
	filled := make([]uint32, len(variance))

	for _, r := range t.Rows {
		for i := range r.Values {
			if r.Values[i].X == nilValue {
				continue
			}

			variance[i] += math.Pow(r.Values[i].X-mean[i], 2)
			filled[i]++
		}
	}

	// 2nd pass based on supplied values
	for i := range variance {
		variance[i] = math.Sqrt(variance[i] / float64(filled[i]))
	}

	return variance
}

func columnSEMs(t *table.Table, stddev []float64) []float64 {
	// check filled values
	sems := make([]float64, len(stddev))
	filled := make([]uint32, len(sems))

	for _, r := range t.Rows {
		for i := range r.Values {
			if r.Values[i].X == nilValue {
				continue
			}

			filled[i]++
		}
	}

	// 2nd pass based on supplied values
	for i := range sems {
		sems[i] = stddev[i] / math.Sqrt(float64(filled[i]))
	}

	return sems
}

func columnCI95s(sem []float64) []float64 {
	// gather values into ci95
	ci95 := make([]float64, len(sem))

	for i := range ci95 {
		ci95[i] = 1.96 * sem[i]
	}

	return ci95
}

func columnCVs(std []float64, mean []float64) []float64 {
	// gather values into cvs
	cvs := make([]float64, len(std))

	for i := range cvs {
		cvs[i] = std[i] / mean[i] * 100.0
	}

	return cvs
}

func RowStats(t *table.Table) {
	t.Keys = append(t.Keys, []string{"MEAN", "STDDEV", "SEM", "CI95", "CV"}...)

	var wg sync.WaitGroup
	for _, row := range t.Rows {
		wg.Add(1)

		go func(r *table.Record) {
			defer wg.Done()

			vals := rowStats(r)
			r.Values = append(r.Values, []table.Value{
				{X: vals[0], UsesUnits: true},
				{X: vals[1], UsesUnits: true},
				{X: vals[2], UsesUnits: true},
				{X: vals[3], UsesUnits: true, Prefix: "±"},
				{X: vals[4], Suffix: "%"},
			}...)
		}(row)
	}
	wg.Wait()
}

func rowStats(r *table.Record) []float64 {
	var sum float64
	var count uint32
	var variance float64

	for i := range r.Values {
		if r.Values[i].X == nilValue {
			continue
		}

		sum += r.Values[i].X
		count++
	}

	mean := sum / float64(count)

	for i := range r.Values {
		if nilValue == r.Values[i].X {
			continue
		}

		variance += math.Pow(r.Values[i].X-mean, 2)
	}

	stddev := math.Sqrt(variance / float64(count))
	sem := stddev / math.Sqrt(float64(count))
	ci95 := 1.96 * sem
	cv := stddev / mean * 100.0

	return []float64{mean, stddev, sem, ci95, cv}
}

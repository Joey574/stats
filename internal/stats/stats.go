package stats

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/Joey574/stats/internal/table"
)

func ColumnStats(t *table.Table) []*table.Record {
	units := t.Rows[0].Units

	avgs := columnAverages(t)
	stds := columnStddevs(t, avgs)
	sems := columnSEMs(t, stds)
	ci95s := columnCI95s(sems)
	cvs := columnCVs(stds, avgs)

	conv := func(x []float64, prefix string, suffix string, usesUnits bool) []table.Value {
		str := make([]table.Value, len(x))

		for i := range x {
			str[i] = table.Value{
				X:         fmt.Sprintf("%.2f", x[i]),
				UsesUnits: usesUnits,
				Prefix:    prefix,
				Suffix:    suffix,
			}
		}
		return str
	}

	avgStr := conv(avgs, "", "", true)
	stdStr := conv(stds, "", "", true)
	semStr := conv(sems, "", "", true)
	ciStr := conv(ci95s, "±", "", true)
	cvStr := conv(cvs, "", "%", false)

	var records []*table.Record
	records = append(records, &table.Record{
		Label:  "MEAN",
		Units:  units,
		Values: avgStr,
	})
	records = append(records, &table.Record{
		Label:  "STDDEV",
		Units:  units,
		Values: stdStr,
	})
	records = append(records, &table.Record{
		Label:  "SEM",
		Units:  units,
		Values: semStr,
	})
	records = append(records, &table.Record{
		Label:  "CI95",
		Units:  units,
		Values: ciStr,
	})
	records = append(records, &table.Record{
		Label:  "CV",
		Units:  "%",
		Values: cvStr,
	})

	return records
}

func columnAverages(t *table.Table) []float64 {
	// gather sums into avgs
	avgs := make([]float64, len(t.Rows[0].Values))
	filled := make([]float64, len(avgs))

	for _, r := range t.Rows {
		for i := range r.Values {
			v, err := strconv.ParseFloat(r.Values[i].X, 64)
			if err != nil {
				continue
			}

			avgs[i] += v
			filled[i]++
		}
	}

	// average out based on supplied values
	for i := range avgs {
		avgs[i] /= filled[i]
	}

	return avgs
}

func columnStddevs(t *table.Table, mean []float64) []float64 {
	// gather sums into variance
	variance := make([]float64, len(mean))
	filled := make([]float64, len(variance))

	for _, r := range t.Rows {
		for i := range r.Values {
			v, err := strconv.ParseFloat(r.Values[i].X, 64)
			if err != nil {
				continue
			}

			variance[i] += math.Pow(v-mean[i], 2)
			filled[i]++
		}
	}

	// 2nd pass based on supplied values
	for i := range variance {
		variance[i] = math.Sqrt(variance[i] / filled[i])
	}

	return variance
}

func columnSEMs(t *table.Table, stddev []float64) []float64 {
	// check filled values
	sems := make([]float64, len(stddev))
	filled := make([]float64, len(sems))

	for _, r := range t.Rows {
		for i := range r.Values {
			_, err := strconv.ParseFloat(r.Values[i].X, 64)
			if err != nil {
				continue
			}

			filled[i]++
		}
	}

	// 2nd pass based on supplied values
	for i := range sems {
		sems[i] = stddev[i] / math.Sqrt(filled[i])
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
				{X: fmt.Sprintf("%.2f", vals[0]), UsesUnits: true},
				{X: fmt.Sprintf("%.2f", vals[1]), UsesUnits: true},
				{X: fmt.Sprintf("%.2f", vals[2]), UsesUnits: true},
				{X: fmt.Sprintf("%.2f", vals[3]), UsesUnits: true, Prefix: "±"},
				{X: fmt.Sprintf("%.2f", vals[4]), Suffix: "%"},
			}...)
		}(row)
	}
	wg.Wait()

}

func rowStats(r *table.Record) []float64 {
	var sum float64
	var count float64
	var variance float64

	for i := range r.Values {
		v, err := strconv.ParseFloat(r.Values[i].X, 64)
		if err != nil {
			continue
		}

		sum += v
		count++
	}

	mean := sum / count

	for i := range r.Values {
		v, err := strconv.ParseFloat(r.Values[i].X, 64)
		if err != nil {
			continue
		}

		variance += math.Pow(v-mean, 2)
	}

	stddev := math.Sqrt(variance / count)
	sem := stddev / math.Sqrt(count)
	ci95 := 1.96 * sem
	cv := stddev / mean * 100.0

	return []float64{mean, stddev, sem, ci95, cv}
}

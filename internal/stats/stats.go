package stats

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Joey574/stats/internal/table"
)

func ColumnStats(t *table.Table) []*table.Record {
	units := t.Rows[0].Units

	avgs := columnAverages(t)
	stds := columnStddevs(t, avgs)
	sems := columnSEMs(t, stds)
	ci95s := columnCI95s(sems)
	cvs := columnCVs(stds, avgs)

	conv := func(x []float64, prefix string, suffix string) []string {
		str := make([]string, len(x))

		for i := range x {
			str[i] = fmt.Sprintf("%s%.2f%s", prefix, x[i], suffix)
		}
		return str
	}

	avgStr := conv(avgs, "", "")
	stdStr := conv(stds, "", "")
	semStr := conv(sems, "", "")
	ciStr := conv(ci95s, "±", "")
	cvStr := conv(cvs, "", "")

	var records []*table.Record
	records = append(records, &table.Record{
		Label:  "mean",
		Units:  units,
		Values: avgStr,
	})
	records = append(records, &table.Record{
		Label:  "stddev",
		Units:  units,
		Values: stdStr,
	})
	records = append(records, &table.Record{
		Label:  "sem",
		Units:  units,
		Values: semStr,
	})
	records = append(records, &table.Record{
		Label:  "ci95",
		Units:  units,
		Values: ciStr,
	})
	records = append(records, &table.Record{
		Label:  "cv",
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
			v, err := strconv.ParseFloat(r.Values[i], 64)
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
			v, err := strconv.ParseFloat(r.Values[i], 64)
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
			_, err := strconv.ParseFloat(r.Values[i], 64)
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
	t.Keys = append(t.Keys, []string{"mean", "stddev", "sem", "ci95", "cv"}...)

	conv := func(x []float64, prefix string, suffix string) []string {
		str := make([]string, len(x))

		for i := range x {
			str[i] = fmt.Sprintf("%s%.2f%s", prefix, x[i], suffix)
		}
		return str
	}

	for _, r := range t.Rows {
		vals := rowStats(r)
		r.Values = append(r.Values, conv(vals, "", "")...)
	}
}

func rowStats(r *table.Record) []float64 {
	var sum float64
	var count float64
	var variance float64

	for i := range r.Values {
		v, err := strconv.ParseFloat(r.Values[i], 64)
		if err != nil {
			continue
		}

		sum += v
		count++
	}

	mean := sum / count

	for i := range r.Values {
		v, err := strconv.ParseFloat(r.Values[i], 64)
		if err != nil {
			continue
		}

		variance += math.Pow(v-mean, 2)
	}

	stddev := math.Sqrt(variance / count)
	sem := stddev / math.Sqrt(count)
	ci95 := 1.96 * sem
	cv := stddev / mean

	return []float64{mean, stddev, sem, ci95, cv}
}

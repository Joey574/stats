package stats

import "math"

type Data struct {
	Stddev float64
	Mean   float64
	Sem    float64
	CI95   float64
	CV     float64
}

type RowData struct {
	Diff  *float64
	Error *float64
}

func CalculateStats(dp []float64) *Data {
	n := float64(len(dp))
	if n == 0 {
		return nil
	}

	// calculate mean
	var sum float64
	for _, v := range dp {
		sum += v
	}
	mean := sum / n

	// calculate variance
	var variance float64
	for _, v := range dp {
		variance += math.Pow(v-mean, 2)
	}
	stddev := math.Sqrt(variance / (n - 1))
	sem := stddev / math.Sqrt(n)
	cv := stddev / mean * 100.0

	// TODO 1.96 is likely too optimistic for my use case, should add lookup table in future
	ci95 := 1.96 * sem

	return &Data{
		Mean:   mean,
		Stddev: stddev,
		Sem:    sem,
		CI95:   ci95,
		CV:     cv,
	}
}

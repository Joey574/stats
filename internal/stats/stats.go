package stats

import (
	"math"

	"github.com/Joey574/stats/internal/table"
)

const nMetrics = 5
const nilValue = math.SmallestNonzeroFloat64

func tableStats(t *table.Table) []float64 {
	rowCount := len(t.Rows)
	rowBlock := make([]float64, nMetrics*rowCount)
	rowFilled := make([]uint32, rowCount)

	// generate statistical information
	for i, r := range t.Rows {

		// collect sums and filled count
		for k := range r.Values {
			if r.Values[k].X == nilValue {
				continue
			}

			rowFilled[i]++
			rowBlock[i*nMetrics] += r.Values[k].X
		}

		// compute mean
		rowBlock[i*nMetrics] /= float64(rowFilled[i])

		// get variance information
		for k := range r.Values {
			if r.Values[k].X == nilValue {
				continue
			}

			rv := r.Values[k].X - rowBlock[i*nMetrics]
			rowBlock[i*nMetrics+1] += rv * rv
		}

		// compute std
		rowBlock[i*nMetrics+1] /= float64(rowFilled[i])

		// compute misc statistic information
		rowBlock[i*nMetrics+2] = rowBlock[i*nMetrics+1] / math.Sqrt(float64(rowFilled[i]))
		rowBlock[i*nMetrics+3] = rowBlock[i*nMetrics+1] / rowBlock[i*nMetrics] * 100.0
		rowBlock[i*nMetrics+4] = 1.96 * rowBlock[i*nMetrics+2]
	}

	return rowBlock
}

package table

import (
	"fmt"
	"strconv"

	"github.com/Knetic/govaluate"
)

func (r *Record) Compose(n int) []string {
	vals := make([]string, n)

	for i := range r.Values {
		if r.Values[i].X == nilValue {
			vals[i] = nilValueRepl
		} else {
			vals[i] = fmt.Sprintf("%s%.3f%s", r.Values[i].Prefix, r.Values[i].X, r.Values[i].Suffix)
		}
	}

	for i := len(r.Values); i < n; i++ {
		vals[i] = nilValueRepl
	}

	return append([]string{r.Label}, vals...)
}

// Modifies record based on the key, val, and provided expr
func (r *Record) Append(key string, val string, expr *govaluate.EvaluableExpression) {
	var ok bool
	var constants string

	switch key {
	case "label":
		r.Label = val
	case "constants":
		// stores row specific constants for applied expressions
		constants = val
	default:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			v = nilValue
		}

		// if expression is nil we just fall through and append the value as is
		if expr != nil {
			params := makeParams(v, constants)
			result, err := expr.Evaluate(params)
			if err != nil {
				v = nilValue
			}

			v, ok = result.(float64)
			if !ok {
				v = nilValue
			}
		}

		r.Values = append(r.Values, Value{
			X: v,
		})
	}
}

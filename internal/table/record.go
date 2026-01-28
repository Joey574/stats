package table

import "fmt"

func (r *Record) Compose(n int) []string {
	vals := make([]string, n)

	for i := range r.Values {
		if r.Values[i].X == nilValue {
			vals[i] = nilValueRepl
		} else {
			vals[i] = fmt.Sprintf("%s%.2f%s", r.Values[i].Prefix, r.Values[i].X, r.Values[i].Suffix)

			if r.Values[i].UsesUnits {
				vals[i] += r.Units
			}
		}
	}

	for i := len(r.Values); i < n; i++ {
		vals[i] = nilValueRepl
	}

	return append([]string{r.Label}, vals...)
}

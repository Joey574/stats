package table

import "fmt"

type Record struct {
	Label string
	Units string

	Values []Value
}

func (r *Record) Compose() []string {
	vals := make([]string, len(r.Values))
	for i := range vals {
		if r.Values[i].X == nilValue {
			vals[i] = nilValueRepl
		} else {
			vals[i] = fmt.Sprintf("%s%.2f%s", r.Values[i].Prefix, r.Values[i].X, r.Values[i].Suffix)

			if r.Values[i].UsesUnits {
				vals[i] += r.Units
			}
		}
	}

	return append([]string{r.Label}, vals...)
}

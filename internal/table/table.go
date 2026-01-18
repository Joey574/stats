package table

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Header struct {
}

type Value struct {
	X         string
	Prefix    string
	Suffix    string
	UsesUnits bool
}

type Table struct {
	Name string
	Keys []string
	Rows []*Record
}

type Record struct {
	Label string
	Units string

	Values []Value
}

const nilTable = "no name"
const nilValue = "-"

func ParseTables(f string) ([]Table, error) {
	bytes, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(strings.NewReader(string(bytes)))

	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	keys := slices.Clone(header)
	keys = slices.DeleteFunc(keys, func(x string) bool { return x == "table" || x == "units" })

	tableMap := make(map[string]*Table)
	for _, row := range records {
		item := Record{}

		var tableName = nilTable
		for i, val := range row {
			name := header[i]

			switch {
			case name == "table":
				tableName = val
			case name == "label":
				item.Label = val
			case name == "units":
				item.Units = val
			default:
				str := val
				if str == "" {
					str = nilValue
				}
				item.Values = append(item.Values, Value{
					X:         val,
					UsesUnits: true,
				})
			}
		}

		if _, ok := tableMap[tableName]; !ok {
			tableMap[tableName] = &Table{
				Name: tableName,
				Keys: keys,
			}
		}

		tableMap[tableName].Rows = append(tableMap[tableName].Rows, &item)
	}

	tables := make([]Table, 0, len(tableMap))
	for _, v := range tableMap {
		tables = append(tables, *v)
	}

	return tables, nil
}

func (t *Table) Headers() []string {
	headers := t.Keys
	headers = slices.DeleteFunc(headers, func(x string) bool {
		return x == "table" || x == "units"
	})

	return headers
}

func (t *Table) Size() (int, int) {
	return len(t.Rows) + 1, len(t.Headers())
}

func (r *Record) Compose() []string {
	vals := make([]string, len(r.Values))
	for i := range vals {
		vals[i] = fmt.Sprintf("%s%s%s", r.Values[i].Prefix, r.Values[i].X, r.Values[i].Suffix)
		if r.Values[i].UsesUnits {
			vals[i] += r.Units
		}
	}

	return append([]string{r.Label}, vals...)
}

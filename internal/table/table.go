package table

import (
	"encoding/csv"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	fixtures "github.com/Joey574/stats/internal/testfixtures"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

type Value struct {
	X         float64
	Prefix    string
	Suffix    string
	UsesUnits bool
}

type Table struct {
	Name string
	Keys []string
	Rows []*Record
}

const nilTable = "no name"
const nilValueRepl = "-"
const nilValue = math.SmallestNonzeroFloat64

func ParseTestTable(f string) (Table, error) {
	file, err := fixtures.TestCSV.Open(f)
	if err != nil {
		return Table{}, err
	}

	reader := csv.NewReader(file)
	reader.ReuseRecord = true

	head, err := reader.Read()
	t, err := ParseTable(reader, slices.Clone(head))
	if err != nil && err != io.EOF {
		return Table{}, err
	}

	return t, nil
}

func ParseTables(f string) ([]Table, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	tables := make([]Table, 0, 1)
	reader := csv.NewReader(file)
	reader.ReuseRecord = true

	for {
		head, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return tables, nil
			}
			return nil, err
		}

		t, err := ParseTable(reader, slices.Clone(head))
		if err != nil && err != io.EOF {
			return nil, err
		}

		tables = append(tables, t)
	}
}

func ParseTable(reader *csv.Reader, header []string) (Table, error) {
	keys := slices.DeleteFunc(slices.Clone(header), func(x string) bool {
		return x == "table" || x == "label" || x == "units"
	})

	name := nilTable
	table := Table{Keys: keys}

	for {
		row, err := reader.Read()
		if err != nil {
			table.Name = name
			return table, err
		}

		item := Record{Values: make([]Value, 0, len(keys))}
		for i, val := range row {
			name := header[i]

			switch name {
			case "table":
				name = val
			case "label":
				item.Label = val
			case "units":
				item.Units = val
			default:
				v, err := strconv.ParseFloat(val, 64)
				if err != nil {
					v = nilValue
				}

				item.Values = append(item.Values, Value{
					X:         v,
					UsesUnits: true,
				})
			}
		}

		table.Rows = append(table.Rows, &item)
	}
}

func (t *Table) Bytes() int64 {
	return int64(8 * len(t.Rows) * len(t.Keys))
}

func (t *Table) Headers() []string {
	return append([]string{"Label"}, t.Keys...)
}

func (c *Table) Dump(renderer tw.Renderer) string {
	if renderer == nil {
		return ""
	}

	var b strings.Builder
	writer := tablewriter.NewTable(&b,
		tablewriter.WithRenderer(renderer))
	writer.Header(c.Headers())

	for _, r := range c.Rows {
		writer.Append(r.Compose(len(c.Keys)))
	}

	writer.Render()
	return b.String()
}

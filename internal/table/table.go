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
	Rows []Record
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
		return slices.Contains(reserved, x)
	})

	name := nilTable
	table := Table{Keys: keys, Rows: make([]Record, 0, 128)}

	for {
		row, err := reader.Read()
		if err != nil {
			table.Name = name
			return table, err
		}

		record := Record{Values: make([]Value, 0, len(keys))}
		for i, val := range row {
			name := header[i]

			switch name {
			case "table":
				name = val
			case "label":
				record.Label = val
			case "units":
				record.Units = val
			default:
				v, err := strconv.ParseFloat(val, 64)
				if err != nil {
					v = nilValue
				}

				record.Values = append(record.Values, Value{
					X:         v,
					UsesUnits: true,
				})
			}
		}

		table.Rows = append(table.Rows, record)
	}
}

func (t *Table) Bytes() int64 {
	return int64(8 * len(t.Rows) * len(t.Keys))
}

func (t *Table) Headers(label string) []string {
	return append([]string{label}, t.Keys...)
}

func (c *Table) Dump(renderer tw.Renderer, label string) string {
	if renderer == nil {
		return ""
	}

	var b strings.Builder
	writer := tablewriter.NewTable(&b,
		tablewriter.WithRenderer(renderer))
	writer.Header(c.Headers(label))

	for _, r := range c.Rows {
		writer.Append(r.Compose(len(c.Keys)))
	}

	writer.Render()
	return b.String()
}

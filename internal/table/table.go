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
	"github.com/Knetic/govaluate"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

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
	t, err := ParseTable(reader, slices.Clone(head), nil)
	if err != nil && err != io.EOF {
		return Table{}, err
	}

	return t, nil
}

func ParseTables(f string, eq string) ([]Table, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	expr, err := govaluate.NewEvaluableExpression(eq)
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

		t, err := ParseTable(reader, slices.Clone(head), expr)
		if err != nil && err != io.EOF {
			return nil, err
		}

		tables = append(tables, t)
	}
}

func ParseTable(reader *csv.Reader, header []string, expr *govaluate.EvaluableExpression) (Table, error) {
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
			case "constants":
				record.Constants = val
			default:
				v, err := strconv.ParseFloat(val, 64)
				if err != nil {
					v = nilValue
				}

				params := makeParams(v, record.Constants)
				result, err := expr.Evaluate(params)
				if err != nil {
					return table, err
				}

				v, ok := result.(float64)
				if !ok {
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

func makeParams(v float64, constants string) map[string]interface{} {
	params := make(map[string]interface{})
	params["x"] = v

	// we now have a slice of strings in the form x=n
	// where x is a string, and n is a numerical value
	consts := strings.SplitSeq(constants, ";")

	for c := range consts {
		s := strings.Split(c, "=")
		if len(s) != 2 {
			continue
		}

		x, err := strconv.ParseFloat(s[1], 64)
		if err != nil {
			continue
		}

		params[s[0]] = x
	}

	return params
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

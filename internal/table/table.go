package table

import (
	"encoding/csv"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

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

	keys := slices.DeleteFunc(slices.Clone(header), func(x string) bool { return x == "table" || x == "units" })

	tableMap := make(map[string]*Table)
	for _, row := range records {
		item := Record{}

		var tableName = nilTable
		for i, val := range row {
			name := header[i]

			switch name {
			case "table":
				tableName = val
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

func (t *Table) Bytes() int64 {
	return int64(8 * len(t.Rows) * len(t.Keys))
}

func (c *Table) Dump(renderer tw.Renderer) string {
	var b strings.Builder
	writer := tablewriter.NewTable(&b,
		tablewriter.WithRenderer(renderer))
	writer.Header(c.Keys)

	for _, r := range c.Rows {
		writer.Append(r.Compose())
	}

	writer.Render()
	return b.String()
}

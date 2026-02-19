package table

import (
	"math"
	"strconv"
	"strings"

	"github.com/Joey574/stats/internal/cli"
	"github.com/olekukonko/tablewriter"
)

const nilValueRepl = "-"
const nilValue = math.SmallestNonzeroFloat64

func NewTable(keys []string) *Table {
	return &Table{
		Keys: keys,
		Rows: make([]Record, 0, 32),
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

func (c *Table) Dump(f cli.CLIArgs) string {
	renderer := f.Renderer()
	if renderer == nil {
		return ""
	}

	var b strings.Builder
	writer := tablewriter.NewTable(&b,
		tablewriter.WithRenderer(renderer))
	writer.Header(c.Headers(f.Label))

	for _, r := range c.Rows {
		writer.Append(r.Compose(len(c.Keys)))
	}

	writer.Render()
	return b.String()
}

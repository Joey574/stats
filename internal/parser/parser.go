package parser

import (
	"bufio"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/Joey574/stats/internal/cli"
	"github.com/Joey574/stats/internal/table"
	fixtures "github.com/Joey574/stats/internal/testfixtures"
	"github.com/Knetic/govaluate"
)

func ParseTestTable(path string) (*table.Table, error) {
	file, err := fixtures.TestCSV.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	t, err := ParseTable(scanner, nil)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return t, nil
}

func ParseTables(f cli.CLIArgs) ([]*table.Table, error) {
	file, err := os.Open(f.File)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// if an expression was passed, create an evaluable expression from it
	var expr *govaluate.EvaluableExpression
	if f.MathEq != "" {
		expr, err = govaluate.NewEvaluableExpression(f.MathEq)
		if err != nil {
			return nil, err
		}
	}

	tables := make([]*table.Table, 0, 1)
	scanner := bufio.NewScanner(file)

	// each table is expected to be delimited by an empty line
	for {
		t, err := ParseTable(scanner, expr)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if t != nil {
			tables = append(tables, t)
		}

		if err == io.EOF {
			return tables, nil
		}
	}
}

func ParseTable(scanner *bufio.Scanner, expr *govaluate.EvaluableExpression) (*table.Table, error) {
	set := false
	var headers []string
	var keys []string
	var t *table.Table

	// parse out the table
	for {
		// we've run out of tokens -> return table and eof
		if !scanner.Scan() {
			return t, io.EOF
		}

		// empty line, means the end of this table -> return table
		txt := scanner.Text()
		if txt == "" {
			return t, nil
		}

		// if the string starts with # ignore the line -> continue
		if strings.TrimSpace(txt)[0] == '#' {
			continue
		}

		row := strings.Split(txt, ",")

		// if we've yet to parse headers, do it now -> continue
		if !set {
			headers = row
			keys := slices.DeleteFunc(slices.Clone(headers), func(x string) bool {
				return slices.Contains(table.Reserved, x)
			})

			t = table.NewTable(keys)
			set = true
			continue
		}

		record := table.Record{Values: make([]table.Value, 0, len(keys))}

		// parse out key, val pairs in the row
		for i, val := range row {
			record.Append(headers[i], val, expr)
		}

		t.Rows = append(t.Rows, record)
	}
}

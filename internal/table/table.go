package table

import (
	"os"

	"github.com/gocarina/gocsv"
)

type Data struct {
	X     *float64 `csv:"x,omitempty"`
	Y     *float64 `csv:"y,omitempty"`
	Truth *float64 `csv:"truth,omitempty"`
	Table string   `csv:"table"`
}

type Row struct {
	X     *float64
	Y     *float64
	Truth *float64
}

type Table struct {
	Name string
	Rows []*Row
}

func ParseTables(f string) ([]*Table, error) {
	in, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	data := []*Data{}
	if err := gocsv.UnmarshalFile(in, &data); err != nil {
		return nil, err
	}

	// read data into their tables
	tables := make(map[string]*Table)
	for _, d := range data {
		if _, ok := tables[d.Table]; !ok {
			tables[d.Table] = &Table{
				Name: d.Table,
			}
		}

		r := &Row{
			X:     d.X,
			Y:     d.Y,
			Truth: d.Truth,
		}
		tables[d.Table].Rows = append(tables[d.Table].Rows, r)
	}

	tableSlice := make([]*Table, 0, len(tables))
	for _, v := range tables {
		tableSlice = append(tableSlice, v)
	}

	return tableSlice, nil
}

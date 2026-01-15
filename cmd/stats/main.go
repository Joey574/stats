package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/Joey574/stats/v2/internal/stats"
	"github.com/Joey574/stats/v2/internal/table"
	"github.com/jessevdk/go-flags"
	"github.com/olekukonko/tablewriter"
)

type CLIArgs struct {
	File string `short:"f" long:"file" description:"csv file to read data from"`
}

type CompiledCol struct {
	Name string
	stats.Data
}

type CompiledRow struct {
	table.Row
	stats.RowData
}

type CompiledTable struct {
	Table *table.Table
	Cols  []*CompiledCol
	Rows  []*CompiledRow
}

func printHeader() {
	fmt.Println(
		`Definitions:
	Mean: The average value among the provided data points
	Stddev: The average distance individual points are from the mean
	SEM (Standard Error of the Mean): The uncertainty in the best estimate
	CI95 (95% Confidence Interval): The range around the mean where we are 95% confident the true value is
	CV (Coefficient of Variation): Meassures how "noisy" the samples are, < 3% is typically considered "good"`)
	fmt.Println()
}

func (t *CompiledTable) Dump() string {
	var b strings.Builder

	// --- Metadata/Column Stats Table ---
	b.WriteString(fmt.Sprintf("Table: %s\n", t.Table.Name))
	colTable := tablewriter.NewWriter(&b)
	colTable.Header([]string{"Key", "Mean", "Stddev", "SEM", "CI95", "CV"})

	for _, k := range t.Cols {
		colTable.Append([]string{
			k.Name,
			fmt.Sprintf("%.2f", k.Mean),
			fmt.Sprintf("%.2f", k.Stddev),
			fmt.Sprintf("%.2f", k.Sem),
			fmt.Sprintf("±%.2f", k.CI95),
			fmt.Sprintf("%.2f%%", k.CV),
		})
	}
	colTable.Render()

	// --- Per Row Table ---
	b.WriteString(fmt.Sprintf("\nTable: %s\n", t.Table.Name))
	rowTable := tablewriter.NewWriter(&b)
	rowTable.Header([]string{"Idx", "Diff", "Error"})
	//rowTable.SetBorder(true)

	for i, r := range t.Rows {
		diffStr, errStr := "N/A", "N/A"
		if r.Diff != nil {
			diffStr = fmt.Sprintf("%.2f%%", *r.Diff)
		}
		if r.Error != nil {
			errStr = fmt.Sprintf("%.2f%%", *r.Error)
		}
		rowTable.Append([]string{fmt.Sprintf("%d", i), diffStr, errStr})
	}
	rowTable.Render()

	return b.String()
}

func main() {
	var f CLIArgs

	// read in arguments
	_, err := flags.Parse(&f)
	if err != nil {
		if !flags.WroteHelp(err) {
			log.Fatalln(err)
		}
	}

	// read in csv
	tables, err := table.ParseTables(f.File)
	if err != nil {
		log.Fatalln(tables)
	}

	// print out compiled data
	printHeader()
	for _, t := range tables {
		ct := compile(t)
		fmt.Println(ct.Dump())
	}
}

func compile(t *table.Table) *CompiledTable {
	ct := &CompiledTable{
		Table: t,
	}

	// generate row statistics
	for _, r := range t.Rows {
		ct.Rows = append(ct.Rows, generateRowStats(r))
	}

	ct.Cols = generateColumnStats(ct.Rows)
	return ct
}

func generateRowStats(r *table.Row) *CompiledRow {
	rs := &CompiledRow{}
	rs.Truth = r.Truth
	rs.X = r.X
	rs.Y = r.Y

	// calculate x->y diff if both values are present
	if r.X != nil && r.Y != nil {
		s := *r.X - *r.Y
		p := *r.X + *r.Y
		val := math.Abs(s / (0.5 * p) * 100.0)
		rs.Diff = &val
	}

	if r.X != nil && r.Truth != nil {
		t := *r.Truth - *r.X
		val := math.Abs(t / *r.Truth * 100.0)
		rs.Error = &val
	}

	return rs
}

func generateColumnStats(rows []*CompiledRow) []*CompiledCol {
	cols := []*CompiledCol{}

	extract := func(metric string) []float64 {
		var vals []float64
		for _, r := range rows {
			switch metric {
			case "X":
				if r.X != nil {
					vals = append(vals, *r.X)
				}
			case "Y":
				if r.Y != nil {
					vals = append(vals, *r.Y)
				}
			case "Diff":
				if r.Diff != nil {
					vals = append(vals, *r.Diff)
				}
			case "Error":
				if r.Error != nil {
					vals = append(vals, *r.Error)
				}
			}
		}
		return vals
	}

	for _, key := range []string{"X", "Y", "Diff", "Error"} {
		dp := extract(key)

		if len(dp) == 0 {
			continue
		}

		cols = append(cols, &CompiledCol{
			Name: key,
			Data: *stats.CalculateStats(dp),
		})
	}

	return cols
}

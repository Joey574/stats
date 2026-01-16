package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/Joey574/stats/internal/stats"
	"github.com/Joey574/stats/internal/table"
	"github.com/jessevdk/go-flags"
)

type CLIArgs struct {
	File string `short:"f" long:"file" description:"csv file to read data from"`
}

const version = "v0.0.2"

func printHeader() {
	fmt.Println(version)
	fmt.Println(
		`
Definitions:
	Mean: The average value among the provided data points
	Stddev: The average distance individual points are from the mean
	SEM (Standard Error of the Mean): The uncertainty in the best estimate
	CI95 (95% Confidence Interval): The range around the mean where we are 95% confident the true value is
	CV (Coefficient of Variation): Meassures how "noisy" the samples are, < 3% is typically considered "good"`)
	fmt.Println()
}

func main() {
	var f CLIArgs

	parser := flags.NewParser(&f, flags.Default)
	parser.Name = "stats"
	parser.ShortDescription = version

	// read in arguments
	_, err := parser.Parse()
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

	printHeader()

	// compile tables in parallel
	var wg sync.WaitGroup
	compiled := make([]*CompiledTable, len(tables))
	for i, t := range tables {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			compiled[idx] = compile(t)
		}(i)
	}
	wg.Wait()

	// dump data
	for _, ct := range compiled {
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

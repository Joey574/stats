package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"

	"github.com/Joey574/stats/internal/stats"
	"github.com/Joey574/stats/internal/table"
	"github.com/jessevdk/go-flags"
)

type CLIArgs struct {
	File    string `short:"f" long:"file" description:"csv file to read data from"`
	Version bool   `long:"version" description:"prints version and exits"`
}

var version string

func printHeader() {
	fmt.Println(version)
	fmt.Println(
		`
Definitions:
	Mean: The average value among the provided data points
	Stddev: The average distance individual points are from the mean
	SEM (Standard Error of the Mean): The uncertainty in the best estimate
	CI95 (95% Confidence Interval): The range around the mean where we are 95% confident the true value is
	CV (Coefficient of Variation): Meassures how "noisy" the samples are, < 3% is typically considered "good

Equations:
	Mean := x̄ = (1/n) Σ x
	Stddev := σ = √ (1/n) Σ (x-x̄)²
	SEM := s = σ / √n
	CI₉₅ := 1.96 * s
	CV := σ / x̄`)
	fmt.Println()
}

func main() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		version = "unknown"
	} else {
		version = info.Main.Version
	}

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
		os.Exit(0)
	}

	if f.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	// read in csv
	tables, err := table.ParseTables(f.File)
	if err != nil {
		log.Fatalln(err)
	}

	printHeader()

	// compile tables in parallel
	var wg sync.WaitGroup
	compiled := make([]stats.CompiledTable, len(tables))
	for i, t := range tables {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			compiled[idx] = stats.CompiledTable{Table: &t}
			compiled[idx].Compile()
		}(i)
	}
	wg.Wait()

	// dump data
	for _, ct := range compiled {
		fmt.Println(ct.Dump())
	}
}

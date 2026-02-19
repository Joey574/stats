package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/Joey574/stats/internal/cli"
	"github.com/Joey574/stats/internal/parser"
	"github.com/Joey574/stats/internal/stats"
)

func main() {
	var f cli.CLIArgs
	f.Parse()

	// read in csv
	tables, err := parser.ParseTables(f)
	if err != nil {
		log.Fatalln(err)
	}

	// compile tables in parallel
	var wg sync.WaitGroup
	compiled := make([]stats.CompiledTable, len(tables))
	for i, t := range tables {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			compiled[idx] = stats.CompiledTable{Table: t}

			if f.FormatTable {
				return
			}

			compiled[idx].CompileDataTable()
		}(i)
	}
	wg.Wait()

	// dump data
	for _, ct := range compiled {
		fmt.Println(ct.Dump(f))
	}
}

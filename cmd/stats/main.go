package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/Joey574/stats/internal/stats"
	"github.com/Joey574/stats/internal/table"
)

var version string

func main() {
	var f CLIArgs
	f.Parse()

	// read in csv
	tables, err := table.ParseTables(f.File)
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
			compiled[idx] = stats.CompiledTable{Table: &t}

			switch {
			case f.ForceTable:
				compiled[idx].CompileForceTable()
			default:
				compiled[idx].CompileDataTable()
			}
		}(i)
	}
	wg.Wait()

	// dump data
	for _, ct := range compiled {
		fmt.Println(ct.Dump(f.Renderer()))
	}
}

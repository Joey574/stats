package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/jessevdk/go-flags"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type RendererArgs struct {
	Text     bool `long:"text"     description:"outputs the table in text format (default: true)"`
	SVG      bool `long:"svg"      description:"outputs the table in svg format"`
	Html     bool `long:"html"     description:"outputs the table in html format"`
	Color    bool `long:"color"    description:"outputs the table in color format"`
	Markdown bool `long:"markdown" description:"outputs the table in markdown format"`
}

type TableArgs struct {
	DataTable  bool `long:"data"  description:"creates a data table (default: true)"`
	ForceTable bool `long:"force" description:"creates and solves a force table"`
}

type CLIArgs struct {
	File    string `short:"f" long:"file"    description:"csv file to read data from"`
	Version bool   `short:"v" long:"version" description:"prints version and exits"`
	Header  bool   `long:"header"`

	RendererArgs `group:"Renderer Args"`
	TableArgs    `group:"Table Args"`
}

const header = `Definitions:
	Mean: The average value among the provided data points
	Stddev: The average distance individual points are from the mean
	SEM (Standard Error of the Mean): The uncertainty in the best estimate
	CI95 (95% Confidence Interval): The range around the mean where we are 95% confident the true value is
	CV (Coefficient of Variation): Measures how "noisy" the samples are, < 3% is typically considered "good"

Equations:
  x̄  Mean   := (1/n) Σ xᵢ
  σ  Stddev := √ (1/n) Σ (xᵢ-x̄)²
  s  SEM    := σ / √n
     CI₉₅   := 1.96 * s
     CV     := σ / x̄`

func (c *CLIArgs) Parse() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		version = "unknown"
	} else {
		version = info.Main.Version
	}

	parser := flags.NewParser(c, flags.Default)
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

	if c.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	if c.Header {
		fmt.Println(header)
		os.Exit(0)
	}
}

func (r *RendererArgs) Renderer() tw.Renderer {
	var ren tw.Renderer
	switch {
	case r.Html:
		ren = renderer.NewHTML()
	case r.SVG:
		ren = renderer.NewSVG()
	case r.Color:
		ren = renderer.NewColorized()
	case r.Markdown:
		ren = renderer.NewMarkdown()
	default:
		ren = renderer.NewBlueprint()
	}

	return ren
}

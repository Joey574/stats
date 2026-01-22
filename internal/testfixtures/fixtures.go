package fixtures

import (
	"embed"
)

//go:embed testdata/csv/*
var TestCSV embed.FS

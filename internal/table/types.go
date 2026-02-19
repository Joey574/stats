package table

type Table struct {
	Name string
	Keys []string
	Rows []Record
}

type Record struct {
	Label     string
	Units     string
	Constants string

	Values []Value
}

type Value struct {
	X         float64
	Sigfigs   int
	Prefix    string
	Suffix    string
	UsesUnits bool
}

package table

type Table struct {
	Keys []string
	Rows []Record
}

type Record struct {
	Label  string
	Values []Value
}

type Value struct {
	X       float64
	Sigfigs int
	Prefix  string
	Suffix  string
}

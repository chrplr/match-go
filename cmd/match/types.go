package main

type Transformation int

const (
	TransNone Transformation = iota
	TransLog10
	TransLength
)

type MissingValueMode int

const (
	MissingMean MissingValueMode = iota
	MissingFixed
)

type Dimension struct {
	Name           string
	File1ID        string
	Col1           int
	File2ID        string
	Col2           int
	Transformation Transformation
	Weight         float64
}

type InputFileConfig struct {
	ID        string
	InputPath string
	OutputPath string
}

type Config struct {
	InputFiles  map[string]*InputFileConfig
	Dimensions  []Dimension
	OutputSize  int
	SummaryFile string
}

type Item struct {
	OriginalValues []string
	Values         []float64 // Transformed and normalized values for each dimension
	ID             int       // Original index in the file
}

type DataSet struct {
	ID    string
	Items []*Item
}

type Tuple struct {
	Items []*Item // Index corresponds to InputFiles order
}

type Solution struct {
	Tuples         []*Tuple
	TotalDistance float64
}

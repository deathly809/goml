package classify

import "github.com/deathly809/goml/data"

// Classifier is a generic interface for a classifier
type Classifier interface {
	Classify([]data.Value) float32
}

// Data is an abstract view of the underlying data
type Data interface {
	Value() []data.Value
	Class() float32
}

// Kernel is a generic Kernel
type Kernel interface {
	Dot([]data.Value, []data.Value) float32
}

package classify

import "github.com/deathly809/gotypes"

// Classifier is a generic interface for a classifier
type Classifier interface {
	Classify([]gotypes.Value) float32
}

// NaiveData is a naive implementation of the
// Data interface
type NaiveData struct {
	Val []gotypes.Value
	Cla float32
}

// Value returns the value array
func (nData *NaiveData) Value() []gotypes.Value {
	return nData.Val
}

// Class returns the classification of this vector
func (nData *NaiveData) Class() float32 {
	return nData.Cla
}

// Data is an abstract view of the underlying data
type Data interface {
	Value() []gotypes.Value
	Class() float32
}

// Kernel is a generic Kernel
type Kernel interface {
	Dot([]gotypes.Value, []gotypes.Value) float32
}

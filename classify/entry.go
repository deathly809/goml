package classify

// Classifier is a generic interface for a classifier
type Classifier interface {
	Classify([]float32) float32
}

// Data is an abstract view of the underlying data
type Data interface {
	Value() []float32
	Class() float32
}

// Kernel is a generic Kernel
type Kernel interface {
	Dot([]float32, []float32) float32
}

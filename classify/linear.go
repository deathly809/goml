package classify

import "github.com/deathly809/gotypes"

type linear struct {
	weights []float64
	bias    float64
}

func (l *linear) Classify(x []gotypes.Value) float32 {
	result := l.bias
	for i, w := range l.weights {
		result += w * x[i].Real()
	}
	return float32(result)
}

// Linear classifier is trained and returned
func Linear(d []Data) Classifier {
	return &linear{}
}

type halfspace struct {
	l *linear
}

func (h *halfspace) Classify(x []gotypes.Value) float32 {
	if h.l.Classify(x) < 0 {
		return -1
	}
	return 1
}

// Halfspace linear predictor
func Halfspace(d []Data) Classifier {
	return &halfspace{
		l: Linear(d).(*linear),
	}
}

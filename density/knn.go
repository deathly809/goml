package density

import (
	"math"

	"github.com/deathly809/goml/data"
	"github.com/deathly809/gotypes"
)

func sqrDiff(x1, x2 gotypes.Value) float64 {
	result := 0.0
	if x1.Type() == gotypes.Array {
		a := x1.Array()
		b := x1.Array()
		for i := range a {
			t := a[i].Real() - b[i].Real()
			result += t * t
		}
	} else {
		result = x1.Real() - x2.Real()
		result *= result
	}
	return result
}

func gaussian(x, m, s gotypes.Value) float64 {
	s2 := s.Real() * s.Real()
	top := math.Exp(-sqrDiff(x, m) / (2.0 * s2))
	bottom := math.Sqrt(math.Pi * 2 * s2)
	return top / bottom
}

type knn struct {
	data data.Set
	h    gotypes.Value
}

func (k *knn) Estimate(data gotypes.Value) float64 {
	result := 0.0
	for i := 0; i < k.data.Rows(); i++ {
		v, _ := gotypes.WrapArray(k.data.Row(i), gotypes.ValueType)
		result += gaussian(v, data, k.h)
	}
	return result / float64(k.data.Rows())
}

// KNN Performs the K-nearest neighbor algorithm
func KNN(input data.Set, K int) Estimator {
	return &knn{data: input}
}

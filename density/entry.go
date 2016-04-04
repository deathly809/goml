package density

import (
	"github.com/deathly809/gotypes"
)

// Estimator will provide a way to estimate the probability
// of a value
type Estimator interface {
	Estimate(data gotypes.Value) float64
}

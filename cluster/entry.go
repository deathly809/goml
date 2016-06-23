package cluster

import "github.com/deathly809/gotypes"

// Metric returns the distance between the values a and b
type Metric func(a []gotypes.Value, b []gotypes.Value) float64

// Point is just a wrapper to confuse map
type Point []float64

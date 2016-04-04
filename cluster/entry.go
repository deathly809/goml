package cluster

import "github.com/deathly809/gotypes"

// Metric returns the distance between the values a and b
type Metric func(a []gotypes.Value, b []float64) float64

/* Empty? */

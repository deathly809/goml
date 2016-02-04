package distributions

import (
	"math"
)

type binom struct {
	n    float64
	prob float64
}

func (b *binom) Pmf(trials []float64) []float64 {
	return []float64{
		math.Pow(b.prob, b.n) * math.Pow(1-b.prob, b.n-trials[0]),
	}
}
func (b *binom) Pdf([]float64) []float64 { return nil }
func (b *binom) Cdf([]float64) []float64 { return nil }
func (b *binom) Mean() []float64         { return nil }
func (b *binom) Median() []float64       { return nil }
func (b *binom) Mode() []float64         { return nil }
func (b *binom) Variance() []float64     { return nil }
func (b *binom) StdDev() []float64       { return nil }

// Binomial creates a binomial distribution
func Binomial(n int, prob float64) {

}

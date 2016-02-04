package distributions

// Distribution represents
type Distribution interface {
	// Discrete distribtions only
	Pmf([]float64) []float64 // Probability mass function

	// Continiouous distributions only
	Pdf([]float64) []float64 // Probability density functions

	// All distributions
	Cdf([]float64) []float64 // Cumulative probability density function
	Mean() []float64         // Mean of the distribution
	Median() []float64       // Median of the distribution
	Mode() []float64         // Mode of the distribution
	Variance() []float64     // Variance of the distribution
	StdDev() []float64       // Standard deviation of the distribution
}

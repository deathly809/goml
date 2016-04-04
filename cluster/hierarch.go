package cluster

// Direction denotes how we build our clusters
type Direction int

/*
 *  Different ways to perform hierarchical clustering
 */
const (
	Divisive      = Direction(iota)
	Agglomerative = Direction(iota)
)

// memory intensive probably
func split(data [][]float64, m Metric, K int) [][]float64 {
	result := make([][]float64, 1)

	return result
}

func build(data [][]float64, m Metric, K int) [][]float64 {
	result := make([][]float64, len(data))

	return result
}

// Hierarchical clustering
func Hierarchical(data [][]float64, m Metric, K int, dir Direction) {
	switch dir {
	case Divisive:
		split(data, m, K)
	case Agglomerative:
		build(data, m, K)
	}
}

package cluster

import "sort"

func arraysEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

//
// Cluster - Holds whatever information is needed to represent a cluster
//
type Cluster struct {
	indices []int
}

// AddPoint - Add a given point to a cluster.  If the
//      the point is nil it is discared.
func (c Cluster) AddPoint(p int) {
	c.indices = append(c.indices, p)
}

// GetPoints - returns the underlying points
//
func (c Cluster) GetPoints() []int {
	return append([]int(nil), c.indices...)
}

// RemovePoint removes a point from the cluster
func (c Cluster) RemovePoint(index int) {
	pos := sort.SearchInts(c.indices, index)
	if pos < len(c.indices) {
		if len(c.indices) == 1 {
			c.indices = nil
		} else {
			c.indices[pos] = c.indices[len(c.indices)-1]
			c.indices = c.indices[:len(c.indices)-1]
		}
	}
}

// Count returns the number of elements in this cluster
func (c Cluster) Count() int {
	return len(c.indices)
}

// CreateCluster creates a cluster from some data
func CreateCluster(data []int) Cluster {
	return Cluster{
		indices: append([]int(nil), data...),
	}
}

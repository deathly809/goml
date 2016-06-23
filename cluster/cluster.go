package cluster

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
	points [][]float64
	n, dim int
}

// AddPoint - Add a given point to a cluster.  If the
//      the point is nil it is discared.
func (c Cluster) AddPoint(p []float64) {
    if p != nil {
        c.points = append(c.points, p)
    }
}

// GetPoints - returns the underlying points
//
func (c Cluster) GetPoints() [][]float64 {
	result := make([][]float64, c.n)
	for i := range c.points {
		result[i] = make([]float64, c.dim)
		copy(result[i], c.points[i])
	}
	return result
}

// RemovePoint removes a point from the cluster
func (c Cluster) RemovePoint(p []float64) {
    // find and remove
    for i := 0 ; i < len(c.points) ; {
        if arraysEqual(p,c.points[i]) {
            c.points[i] = c.points[c.n - 1]
            c.n--
        }else {
            i++
        }
    }

    // shrink if needed
    c.points = c.points[:c.n]
}

// Count returns the number of elements in this cluster
func (c Cluster) Count() int {
	return c.n
}

package cluster

import "sort"

// TODO: Try to incorporate some type of kd-tree datastructure to reduce runtime to ~O(n log(n))

// This is based off of apache commons stuff.  Can't remember name right now

type pointStatus int

/*
 *  ENUM for SEEN and UNSEEN
 */
const (
	_Unseen  = pointStatus(iota)
	_Cluster = pointStatus(iota)
	_Noise   = pointStatus(iota)
)

func withinEpsilon(a, b Point, eps float64) bool {
	sum := 0.0
	for i := 0; i < len(a); i++ {
		t := a[i] * b[i]
		sum += t * t
		if sum > eps*eps {
			return false
		}
	}
	return sum <= eps*eps
}

func getNeighbors(data []Point, start Point, index int, epsilon float64) []int {
	result := []int{index}
	for i, d := range data {
		if i != index {
			if withinEpsilon(d, start, epsilon) {
				result = append(result, i)
			}
		}
	}
	sort.Ints(result)
	return result
}

func merge(a, b []int) []int {
	result := []int(nil)

	for len(a) > 0 && len(b) > 0 {
		if a[0] < b[0] {
			result = append(result, a[0])
			a = a[1:]
		} else {
			result = append(result, b[0])
			b = b[1:]
		}
	}

	result = append(result, a...)
	result = append(result, b...)

	return result
}

func expandCluster(data []Point, cluster []int, seen []pointStatus, epsilon float64) []int {
	tmp := []int(nil)
	for i := range data {
		if seen[i] == _Cluster {
			if sort.SearchInts(cluster, i) == len(cluster) || sort.SearchInts(tmp, i) == len(tmp) {
				for _, j := range cluster {
					if withinEpsilon(data[i], data[j], epsilon) {
						tmp = append(tmp, i)
						seen[i] = _Cluster
					} else {
						seen[i] = _Noise
					}
				}
			}
		}
	}
	return merge(cluster, tmp)
}

// DBSCAN performs density based clustering.
func DBSCAN(data []Point, k int, epsilon float64) []Cluster {

	clusters := []Cluster(nil)
	seen := make([]pointStatus, len(data))
	for i, d := range data {
		if seen[i] == _Unseen {
			neighbors := getNeighbors(data, d, i, epsilon)
			if len(neighbors) >= k {
				seen[i] = _Cluster
				newCluster := CreateCluster(expandCluster(data, neighbors, seen, epsilon))
				clusters = append(clusters, newCluster)
			} else {
				seen[i] = _Noise
			}
		}
	}
	return clusters
}

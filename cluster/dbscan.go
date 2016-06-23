package cluster

// This is based off of apache commons stuff.  Can't remember name right now

type pointStatus int

/*
 *  ENUM for SEEN and UNSEEN
 */
const (
	CLUSTER = pointStatus(iota)
	NOISE   = pointStatus(iota)
)

func getNeighbors(data []Point, start Point, epsilon float64) []Point {
    return nil
}

// DBSCAN performs density based clustering.
func DBSCAN(data []Point, k int, epsilon float64) []Cluster {

	clusters := []Cluster(nil)
	seen := make(map[Point]pointStatus)
	for _, d := range data {
		if v, ok := seen[d]; !ok {
            neighbors := getNeighbors(data,d,epsilon);
            if(len(neighbors) >= k) {
                seen[d] = CLUSTER;
            }else {
                seen[d] = NOISE
		    }
        }
	}
	return clusters
}

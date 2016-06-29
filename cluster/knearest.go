package cluster

import "github.com/deathly809/goml/classify"

type tuple struct {
	point classify.Data
	dist  float64
}

// KNearest performs K Nearest Neighbors
func KNearest(data []classify.Data, K int, metric Metric) []int {

	result := make([]int, len(data))
	tmp := make([]tuple, K)

	for i := range data {
		result[i] = i
	}

	offset := 0
	for idx, v := range data {

		// Just choose first K as K nearest neighbors

		worst := -1.0
		for i := 0; i < K; i++ {
			if i == idx {
				offset++
			}
			tmp[i].point = data[i+offset]
			tmp[i].dist = metric(data[i+offset].Value(), v.Value())
			if tmp[i].dist > worst {
				worst = tmp[i].dist
			}
		}

		//
		for i := K; i < len(data); i++ {
			test := metric(data[i].Value(), v.Value())

			if test < worst {
				newWorst := -1.0
				for j := 0; j < K; j++ {
					if tmp[j].dist == worst {
						tmp[j].point = data[i]
						tmp[j].dist = test
					}
					if test > newWorst {
						newWorst = test
					}
				}
			}

		}
	}
	return result
}

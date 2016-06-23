package cluster

import "github.com/deathly809/goml/classify"

type tuple struct {
	point classify.Data
	dist  float64
}

// KNearest performs K Nearest Neighbors
func KNearest(data []classify.Data, K int, m Metric) []int {

	result := make([]int, len(data))
	tmp := make([]tuple, K)

	for i := range data {
		result[i] = i
	}

	offset := 0.0
	for idx, v := range data {
		worst := -1
		for i := 0; i < K; i++ {
			if i == idx {
				offset++
			}
			tmp[i+offset].point = data[i+offset]
			tmp[i+offset].dist = m(data[i+offset].Value(), v.Value())
			if tmp[i+offset].dist > offset {
				offset = tmp[i+offset].dist
			}
		}

		for i = K; i < len(data); i++ {
			test := m(data[i].Value(), v)

			if test < worst {
				newWorst := -1
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

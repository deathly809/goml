package cluster

import (
	"math"
	"math/rand"
	"sync"

	"github.com/deathly809/gods/queue"
	"github.com/deathly809/gomath"
	"github.com/deathly809/goml/classify"
	"github.com/deathly809/parallels"
)

func average(q queue.Queue, features int) []float64 {
	mean := make([]float64, features)
	for q.Count() > 0 {
		fv := q.Dequeue().([]float64)
		for i, v := range fv {
			mean[i] += v
		}
	}
	for i := range mean {
		mean[i] = (math.Sqrt(mean[i]) / float64(features))
	}
	return mean
}

func same(a, b []float64) bool {
	val := 0.0
	for i := range a {
		t := a[i] - b[i]
		val += t * t
	}
	val = math.Sqrt(val)
	return val < 1E-5
}

// KMeans clustering
// TODO: Change to not only be numeric
func KMeans(data []classify.Data, K int, m Metric, means [][]float64) [][]float64 {

	features := len(data[0].Value())

	if means == nil {
		/* Todo: randomly initialize k means */
		means = make([][]float64, K)
		for i := range means {
			means[i] = make([]float64, features)
			val := data[rand.Intn(len(data))]
			for j, v := range val.Value() {
				means[i][j] = v.Real()
			}
		}
	}

	counts := make([]int, K)
	newMeans := make([][]float64, K)
	for i := range newMeans {
		newMeans[i] = make([]float64, features)
	}

	maxElemPerThread := 10000 // 1K
	iterations := gomath.MaxInt(1, (len(data)+maxElemPerThread-1)/maxElemPerThread)
	elemPerThread := (len(data) + iterations - 1) / iterations

	change := true
	for change {
		change = false
		var lock sync.Mutex
		parallels.Foreach(func(i int) {
			// List of stuff
			nMs := make([][]float64, K)
			cs := make([]int, K)

			for i := range nMs {
				nMs[i] = make([]float64, features)
			}

			start := i * elemPerThread
			end := gomath.MinInt(start+elemPerThread, len(data))

			for start < end {
				d := data[start]
				closest := 0
				dist := m(d.Value(), means[0])
				for k := 1; k < K; k++ {
					nDist := m(d.Value(), means[k])
					if nDist < dist {
						dist = nDist
						closest = k
					}
				}

				// Save
				cs[closest]++
				for i := 0; i < features; i++ {
					t := d.Value()[i].Real()
					nMs[closest][i] += t
				}

				start++
			}

			lock.Lock()
			defer lock.Unlock()

			for i := 0; i < K; i++ {
				for j := 0; j < features; j++ {
					newMeans[i][j] += nMs[i][j]
				}
				counts[i] += cs[i]
			}
		}, iterations)

		// Update means
		for i := range newMeans {
			for j := range newMeans[i] {
				newMeans[i][j] /= float64(counts[i])
			}
			if !same(newMeans[i], means[i]) {
				change = true
			}
			for j := range newMeans[i] {
				means[i][j] = newMeans[i][j]
				newMeans[i][j] = 0
			}
			counts[i] = 0
		}
	}

	return means
}

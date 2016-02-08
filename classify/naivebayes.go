package classify

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/deathly809/gomath"
	"github.com/deathly809/goml/data"
)

/*
 * Reminder : The Naive Bayes
 *
 *  max_C_i [P(C_i,X)]
 *
 *  P(C|X) = (P(X|C))
 *
 */

type gaussian struct {
	mean     float64
	stddev   float64
	constant float64
}

func (g *gaussian) Evaluate(x float64) float64 {
	exp := (x - g.mean) / g.stddev
	return g.constant * math.Exp(-0.5*exp*exp)
}

type naive struct {
	pGaussians []gaussian
	nGaussians []gaussian
	pPositive  float64
	pNegative  float64
}

func (n *naive) Classify(v []data.Value) float32 {
	result := 0.0
	pPositive := n.pPositive
	pNegative := n.pNegative

	for col := 0; col < len(v); col++ {
		pPositive *= gomath.MaxFloat64(0.001, n.pGaussians[col].Evaluate(v[col].Real()))
		pNegative *= gomath.MaxFloat64(0.001, n.nGaussians[col].Evaluate(v[col].Real()))
	}
	if pPositive > pNegative {
		result = +1
	} else if pPositive < pNegative {
		result = -1
	}
	return float32(result)
}

func clean(s string) string {
	return strings.TrimSpace(s)
}

func handleColumn(data []Data, col int) (gaussian, gaussian) {
	var pCount, nCount float64
	var pSum, nSum float64

	pProb := make(map[float64]float64)
	nProb := make(map[float64]float64)

	// Construct means
	for _, r := range data {
		v := r.Value()[col].Real()
		c := r.Class()
		if c < 0 {
			nSum += v
			nCount++
			prob := nProb[v]
			nProb[v] = prob + 1
		} else if c > 0 {
			pSum += v
			pCount++
			prob := pProb[v]
			pProb[v] = prob + 1
		}
	}

	for k, v := range pProb {
		pProb[k] = v / pCount
	}

	for k, v := range nProb {
		nProb[k] = v / pCount
	}

	pMean := pSum / pCount
	nMean := nSum / nCount

	// Construct standard deviation
	pStdDev := 0.0
	nStdDev := 0.0
	for _, r := range data {
		v := r.Value()[col].Real()
		c := r.Class()
		if c < 0 {
			nStdDev = nProb[v] * (v - nMean) * (v - nMean)
		} else if c > 0 {
			pStdDev = pProb[v] * (v - pMean) * (v - pMean)
		}
	}

	pStdDev = math.Sqrt(pStdDev)
	nStdDev = math.Sqrt(nStdDev)

	pStdDev = gomath.ClampFloat64(0.001, pStdDev, pStdDev)
	nStdDev = gomath.ClampFloat64(0.001, nStdDev, nStdDev)

	constant := math.Sqrt2 * math.SqrtPi
	pGaussian := gaussian{
		mean:     pMean,
		stddev:   pStdDev,
		constant: constant * pStdDev,
	}

	nGaussian := gaussian{
		mean:     nMean,
		stddev:   nStdDev,
		constant: constant * nStdDev,
	}
	return pGaussian, nGaussian
}

// New constructs a new naive bayes classifier
func New(data []Data) Classifier {
	result := &naive{}
	cols := len(data[0].Value())

	result.pGaussians = make([]gaussian, cols)
	result.nGaussians = make([]gaussian, cols)

	for _, r := range data {
		if r.Class() > 0 {
			result.pPositive++
		} else if r.Class() < 0 {
			result.pNegative++
		} else {
			log.Fatal("Classification cannot be 0.0")
		}
	}

	fmt.Printf("Number positive=%0.0f, Number of negative=%0.0f\n", result.pPositive, result.pNegative)

	result.pPositive /= float64(len(data))
	result.pNegative /= float64(len(data))

	fmt.Printf("Number positive=%0.2f, Number of negative=%0.2f\n", result.pPositive, result.pNegative)

	for col := 0; col < cols; col++ {
		result.pGaussians[col], result.nGaussians[col] = handleColumn(data, col)
	}

	return result
}

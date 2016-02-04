package classify

import (
	"strings"

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

type naive struct {
	pCount    []map[string]int
	nCount    []map[string]int
	pPositive float32
	pNegative float32
	n         float32
}

func (n *naive) Classify(v []data.Value) float32 {
	result := float32(1.0)
	// Find minimum class -1/+1 (<0,>0)
	pPositive := result
	pNegative := result

	for col := 0; col < len(v); col++ {
		txt := clean(v[col].Text())

		if count, ok := n.pCount[col][txt]; ok {
			pPositive *= float32(count)
		}

		if count, ok := n.nCount[col][txt]; ok {
			pNegative *= float32(count)
		}
	}

	pPositive /= n.pPositive
	pNegative /= n.pNegative

	if pPositive > pNegative {
		result = pPositive
	} else {
		result = pNegative
	}
	return result
}

// Only need the counts of each distinct value/column
// Note: Issues with real values?

func clean(s string) string {
	return strings.TrimSpace(s)
}

func handleColumn(data []Data, col int) (map[string]int, map[string]int) {
	pResult := make(map[string]int)
	nResult := make(map[string]int)

	for _, r := range data {
		txt := clean(r.Value()[col].Text())
		if r.Class() > 0 {
			val, _ := pResult[txt]
			pResult[txt] = val + 1
		} else if r.Class() < 0 {
			val, _ := nResult[txt]
			nResult[txt] = val + 1
		}
	}
	return pResult, nResult
}

// New constructs a new naive bayes classifier
func New(data []Data) Classifier {
	result := &naive{}
	cols := len(data[0].Value())

	result.pCount = make([]map[string]int, cols)
	result.nCount = make([]map[string]int, cols)

	for _, r := range data {
		if r.Class() > 0 {
			result.pPositive++
		} else if r.Class() < 0 {
			result.pNegative++
		}
	}

	for col := 0; col < cols; col++ {
		result.pCount[col], result.nCount[col] = handleColumn(data, col)
	}
	return result
}

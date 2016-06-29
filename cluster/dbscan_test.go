package cluster

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

func loadData(filename string) []Point {
	f, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	input := bufio.NewScanner(f)

	// How many
	input.Scan()
	line := input.Text()

	// columns?
	input.Scan()
	line = strings.Split(input.Text(), " ")[1]
	tmpInt64, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		panic(err)
	}
	dim := int(tmpInt64)

	// uh...
	input.Scan()
	input.Text()

	// Labels
	input.Scan()
	input.Text()

	result := []Point(nil)

	// Data
	for input.Scan() {
		line := strings.Split(input.Text(), "	")
		p := []float64(nil)
		for i := 1; i < dim; i++ {
			fl, err := strconv.ParseFloat(line[i], 64)
			if err != nil {
				panic(err)
			}
			p = append(p, fl)
		}
		result = append(result, p)
	}
	return result
}

func TestDBSCANHepta(t *testing.T) {
	ExpectedClusters := 7
	k := 10
	epsilon := 0.5
	data := loadData("../datasets/Hepta.lrn")

	clusters := DBSCAN(data, k, epsilon)

	if len(clusters) != ExpectedClusters {
		t.Log(len(clusters))
	}
}

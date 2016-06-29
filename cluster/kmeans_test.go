package cluster

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/deathly809/gomath"
	"github.com/deathly809/goml/classify"
	"github.com/deathly809/gotypes"
)

type mydata struct {
	value []gotypes.Value
	class float32
}

func (m *mydata) Value() []gotypes.Value {
	return m.value
}

func (m *mydata) Class() float32 {
	return m.class
}

func met(a []gotypes.Value, b []gotypes.Value) float64 {
	result := 0.0
	for i := range a {
		t := a[i].Real() - b[i].Real()
		result += t * t
	}
	return result
}

var (
	LargeK    = 1000000 // 10M
	LargeData []classify.Data
)

func init() {
	features := 8
	jump := 1000
	var wg sync.WaitGroup
	LargeData = make([]classify.Data, LargeK)
	for i := 0; i < LargeK; i += jump {
		wg.Add(1)
		go func(i int) {
			start := i
			end := gomath.MinInt(LargeK, i+jump)
			for start < end {
				vArray := make([]gotypes.Value, features)
				LargeData[start] = &mydata{value: vArray, class: 0}
				for k := range vArray {
					vArray[k] = gotypes.WrapReal(rand.Float64())
				}
				start++
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func wrapper(data []classify.Data, K int, met Metric, t *testing.T) {
	means := KMeans(data, K, met, nil)

	if means == nil {
		t.FailNow()
	}

	for _, m := range means {
		fmt.Println(m)
	}
}

func createValue(text string) gotypes.Value {
	val, _ := strconv.ParseFloat(text, 64)
	return gotypes.WrapReal(val)
}

func createData(csv string, class float32) classify.Data {
	split := strings.Split(csv, ",")
	vArray := []gotypes.Value(nil)
	for _, s := range split {
		vArray = append(vArray, createValue(s))
	}
	return &mydata{
		value: vArray,
		class: class,
	}
}

func TestCreateKMean(t *testing.T) {
	K := 2
	data := []classify.Data{
		createData("1.0, 2.0", 0),
		createData("2.0, 3.0", 0),
		createData("2.0, 2.0", 0),
		createData("4.0, 3.0", 0),
		createData("3.0, 2.0", 0),
		createData("-1.0, 2.0", 0),
		createData("-2.0, 3.0", 0),
		createData("-2.0, 2.0", 0),
		createData("-4.0, 4.0", 0),
		createData("-3.0, 2.0", 0),
	}
	wrapper(data, K, met, t)
}

func TestCreateLarge(t *testing.T) {
	wrapper(LargeData, 2, met, t)
}

func BenchmarkLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		means := KMeans(LargeData, 10, met, nil)
		if means == nil {
			fmt.Println("error")
		}
	}
}

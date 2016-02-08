package classify

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/deathly809/goml/data"
)

type mydata struct {
	value []data.Value
	class float32
}

func (m *mydata) Value() []data.Value {
	return m.value
}

func (m *mydata) Class() float32 {
	return m.class
}

type myvalue struct {
	real    float64
	theType data.Type
}

func (m *myvalue) Text() string {
	return strconv.FormatFloat(m.real, 'E', -1, 64)
}

func (m *myvalue) Integer() int64 {
	return 0
}

func (m *myvalue) Real() float64 {
	return m.real
}

func (m *myvalue) Boolean() bool {
	return true
}

func (m *myvalue) Type() data.Type {
	return m.theType
}

func (m *myvalue) Initialized() bool {
	return true
}

func createValue(text string) data.Value {
	val, _ := strconv.ParseFloat(text, 64)
	return &myvalue{
		real:    val,
		theType: data.Real,
	}
}

func createData(csv string, class float32) Data {
	split := strings.Split(csv, ",")
	vArray := []data.Value(nil)
	for _, s := range split {
		vArray = append(vArray, createValue(s))
	}
	return &mydata{
		value: vArray,
		class: class,
	}
}

func TestNaiveBayes(t *testing.T) {
	input := []Data{
		createData("6,180,12", 1),
		createData("5.92,190,11", 1),
		createData("5.58,170,12", 1),
		createData("5.92,165,10", 1),
		createData("5,100,6", -1),
		createData("5.5,150,8", -1),
		createData("5.42,130,7", -1),
		createData("5.75,150,9", -1),
	}
	classifier := New(input)
	for _, r := range input {
		fmt.Println(classifier.Classify(r.Value()))
	}
}

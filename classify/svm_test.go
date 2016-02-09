package classify

import (
	"testing"

	"github.com/deathly809/gotypes"
)

type herpKern struct{}

func (herp *herpKern) Dot(a, b []gotypes.Value) float32 {
	return 0.0
}

func TestSVM(t *testing.T) {

	data := []Data(nil)
	for i := 0.0; i < 5.0; i += 0.1 {
		data = append(
			data,
			&mydata{
				value: []gotypes.Value{
					gotypes.WrapReal(i),
				},
				class: -1.0,
			},
		)
	}
	for i := 5.0; i <= 10.0; i += 0.1 {
		data = append(
			data,
			&mydata{
				value: []gotypes.Value{
					gotypes.WrapReal(i),
				},
				class: 1.0,
			},
		)
	}
	Alpha := float32(0.5)
	C := float32(0.5)

	svm := NewSVMClassifer(data, &herpKern{}, Alpha, C)
	if svm == nil {
		t.FailNow()
	}
}

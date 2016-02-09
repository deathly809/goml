package data

import (
	"strings"
	"testing"
)

func TestLoadCSVNoHeader(t *testing.T) {
	data := "Jeffrey,Robinson,32\nGabriel,Loewen,31"
	const (
		ExpectedNumRows = 2
		ExpectedNumCols = 3
	)

	dataSet, err := LoadCSV(strings.NewReader(data), false)
	if err != nil {
		t.Logf("Error: %s", err.Error())
		t.FailNow()
	}
	if dataSet == nil {
		t.Logf("dataSet is nil")
		t.FailNow()
	}

	if dataSet.Rows() != ExpectedNumRows {
		t.Logf("Expected number of rows = %d, Number found = %d\n", ExpectedNumRows, dataSet.Rows())
		t.FailNow()
	}

	if dataSet.Columns() != ExpectedNumCols {
		t.Logf("Expected number of rows = %d, Number found = %d\n", ExpectedNumCols, dataSet.Columns())
		t.FailNow()
	}

}

func TestLoadCSVWithHeader(t *testing.T) {
	data := "First,Last,Age\nJeffrey,Robinson,32\nGabriel,Loewen,31"
	const (
		ExpectedNumRows = 2
		ExpectedNumCols = 3
	)

	dataSet, err := LoadCSV(strings.NewReader(data), true)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if dataSet == nil {
		t.Logf("dataSet is nil")
		t.FailNow()
	}

	if dataSet.Rows() != ExpectedNumRows {
		t.Logf("Expected number of rows = %d, Number found = %d\n", ExpectedNumRows, dataSet.Rows())
		t.FailNow()
	}

	if dataSet.Columns() != ExpectedNumCols {
		t.Logf("Expected number of rows = %d, Number found = %d\n", ExpectedNumCols, dataSet.Columns())
		t.FailNow()
	}

	if dataSet.Header() == nil {
		t.Log("Expected a header but none found")
		t.FailNow()
	}

	if len(dataSet.Header()) != ExpectedNumCols {
		t.Logf("Expected a header with %d columns, found %d\n", ExpectedNumCols, len(dataSet.Header()))
		t.FailNow()
	}
}

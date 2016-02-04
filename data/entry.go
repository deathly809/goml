package data

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

// Type is the underlying type of the data
type Type int

const (
	// Unknown means we don't know what the underlying data type is yet
	Unknown = Type(iota)
	// Real is a floating point value
	Real = Type(iota)
	// Integer is an integral value
	Integer = Type(iota)
	// Boolean is a true/false
	Boolean = Type(iota)
	// Text means we could not put it in any of the others
	Text = Type(iota)
)

// Value is a wrapper for data
type Value interface {
	Type() Type
	Integer() int64
	Real() float64
	Boolean() bool
	Text() string

	Initialized() bool
}

// Set is the interface for a
// tabular set of data
type Set interface {
	Row(int) []Value

	Header() []string
	ColumnTypes() []Type

	Rows() int
	Columns() int
}

func parseRow(row []string, types []Type) []Value {
	result := make([]Value, len(row))
	for i, d := range row {
		v := &value{
			dataType: types[i],
			text:     d,
		}

		var err error

		if types[i] != Unknown {
			switch types[i] {
			case Real:
				v.real, err = strconv.ParseFloat(strings.TrimSpace(d), 64)
			case Integer:
				v.integer, _ = strconv.ParseInt(strings.TrimSpace(d), 10, 64)
			case Boolean:
				v.boolean, _ = strconv.ParseBool(strings.TrimSpace(d))
			}
			v.initialized = (err == nil)
		} else {
			v.initialized = (d != "")
			if v.initialized {
				trimmed := strings.TrimSpace(d)
				types[i], v.dataType = Text, Text
				if v.real, err = strconv.ParseFloat(trimmed, 64); err == nil {
					types[i], v.dataType = Real, Real
				} else if v.integer, err = strconv.ParseInt(trimmed, 10, 64); err == nil {
					types[i], v.dataType = Integer, Integer
				} else if v.boolean, err = strconv.ParseBool(trimmed); err == nil {
					types[i], v.dataType = Boolean, Boolean
				}
			}

		}
		result[i] = v
	}
	return result
}

// LoadCSV will load a data set from a CSV
// source
func LoadCSV(input io.Reader, header bool) Set {
	result := &set{}
	reader := csv.NewReader(input)

	reader.TrimLeadingSpace = true

	if header {
		top, err := reader.Read()
		if err != nil {
			panic("Expected a header, none found")
		}
		for i := range top {
			top[i] = strings.TrimSpace(top[i])
		}
		result.header = top
	}
	for row, err := reader.Read(); err == nil; row, err = reader.Read() {
		if result.cTypes == nil {
			result.cols = len(row)
			result.cTypes = make([]Type, len(row))
		}
		result.data = append(result.data, parseRow(row, result.cTypes))
	}
	result.rows = len(result.data)

	return result
}

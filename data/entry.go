package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/deathly809/gotypes"
)

// Set is the interface for a
// tabular set of data
type Set interface {
	Row(int) []gotypes.Value

	Header() []string
	ColumnTypes() []gotypes.Type

	Rows() int
	Columns() int
}

func parseColumn(text string, theType gotypes.Type) (result gotypes.Value, retType gotypes.Type, err error) {
	if theType == gotypes.Unknown {
		if result, retType, err = parseColumn(text, gotypes.Real); err == nil {
			/* Empty */
		} else if result, retType, err = parseColumn(text, gotypes.Integer); err == nil {
			/* Empty */
		} else if result, retType, err = parseColumn(text, gotypes.Boolean); err == nil {
			/* Empty */
		} else if result, retType, err = parseColumn(text, gotypes.Text); err == nil {
			/* Empty */
		} else {
			panic("Something went seriously wrong")
		}
	} else {

		var (
			real    float64
			integer int64
			boolean bool
		)
		retType = theType

		switch theType {
		case gotypes.Real:
			if real, err = strconv.ParseFloat(strings.TrimSpace(text), 64); err != nil {
				break
			} else {
				result = gotypes.WrapReal(real)
				if result == nil {
					panic("real...")
				}
			}
		case gotypes.Integer:
			if integer, err = strconv.ParseInt(strings.TrimSpace(text), 10, 64); err != nil {
				break
			} else {
				result = gotypes.WrapInteger(integer)
				if result == nil {
					panic("integer...")
				}
			}
		case gotypes.Boolean:
			if boolean, err = strconv.ParseBool(strings.TrimSpace(text)); err != nil {
				break
			} else {
				result = gotypes.WrapBoolean(boolean)
				if result == nil {
					panic("boolean...")
				}
			}
		case gotypes.Text:
			result = gotypes.WrapText(text)
			if result == nil {
				panic("text...")
			}
		default:
			panic(fmt.Sprintf("%s: %d\n", "Unknown type", theType))
		}

		// In the else
		if result == nil && err == nil {
			panic(fmt.Sprintf("%s: %d", "no one caught...", theType))
		}
	}
	return
}

func parseRow(row []string, types []gotypes.Type) (result []gotypes.Value, err error) {
	result = make([]gotypes.Value, len(row))
	for i, d := range row {
		if result[i], types[i], err = parseColumn(d, types[i]); err != nil {
			panic(err.Error())
		}
		if result[i] == nil {
			panic(fmt.Sprintf("%s: %d", "nil result", types[i]))
		}
	}
	return
}

// LoadCSV will load a data set from a CSV
// source
func LoadCSV(input io.Reader, header bool) (r Set, err error) {
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
			result.cTypes = make([]gotypes.Type, len(row))
		}
		r, err := parseRow(row, result.cTypes)
		if err != nil {
			result = nil
			break
		}
		result.data = append(result.data, r)
		result.rows++
	}
	r = result
	return
}

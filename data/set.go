package data

import "github.com/deathly809/gotypes"

type set struct {
	cTypes []gotypes.Type
	data   [][]gotypes.Value
	rows   int
	cols   int

	header []string
}

func (s *set) ColumnTypes() []gotypes.Type {
	return s.cTypes
}

func (s *set) Row(row int) []gotypes.Value {
	result := []gotypes.Value(nil)
	if row >= 0 && row < s.rows {
		result = s.data[row]
	}
	return result
}

func (s *set) Rows() int {
	return s.rows
}

func (s *set) Columns() int {
	return s.cols
}

func (s *set) Header() []string {
	return s.header
}

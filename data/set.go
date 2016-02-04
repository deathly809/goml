package data

type set struct {
	cTypes []Type
	data   [][]Value
	rows   int
	cols   int

	header []string
}

func (s *set) ColumnTypes() []Type {
	return s.cTypes
}

func (s *set) Row(row int) []Value {
	result := []Value(nil)
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

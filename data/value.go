package data

type value struct {
	dataType Type
	integer  int64
	real     float64
	boolean  bool
	text     string

	initialized bool
}

func (v *value) Type() Type {
	return v.dataType
}

func (v *value) Integer() int64 {
	return v.integer
}

func (v *value) Real() float64 {
	return v.real
}

func (v *value) Boolean() bool {
	return v.boolean
}

func (v *value) Text() string {
	return v.text
}

func (v *value) String() string {
	return v.text
}

func (v *value) Initialized() bool {
	return v.initialized
}

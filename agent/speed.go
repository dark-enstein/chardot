package agent

type Speed int

func (s Speed) Ptr() *Speed {
	return &s
}

func (s Speed) Int() Coordinate {
	return Coordinate(s)
}

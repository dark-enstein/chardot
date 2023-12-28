package agent

import (
	"fmt"
	"github.com/dark-enstein/chardot/internal/ilog"
	"strconv"
)

// Coordinate wraps a standard int and adds functionalities for manipulating coordinates.
type Coordinate int

// Int converts a Coordinate to its underlying integer value.
func (p *Coordinate) Int() int {
	return int(*p)
}

// assume sets the Coordinate to the provided integer value.
func (p *Coordinate) assume(i int) {
	*p = Coordinate(i)
}

// add increments the Coordinate by the provided integer value.
func (p *Coordinate) add(i int) {
	po := Coordinate(i)
	*p += po
}

// Sign determines the Sign of the Coordinate, returning POS (true) or NEG (false).
func (p *Coordinate) Sign() Sign {
	if strconv.FormatInt(int64(*p), 10)[0] != '-' {
		return POS
	}
	return NEG
}

// Negate turns a positive Coordinate integer into a negative one
func (p *Coordinate) Negate() error {
	i, err := strconv.ParseInt(fmt.Sprintf("-%v", p), 10, 0)
	if err != nil {
		return ilog.CheckErrAll(err, "error occured while negating")()
	}
	p.assume(int(i))
	return nil
}

// MustNegate turns a positive Coordinate integer into a negative one
func (p *Coordinate) MustNegate() {
	if err := p.Negate(); err != nil {
		panic(err)
	}
}

// Denegate turns a negative Coordinate integer into a positive one
func (p *Coordinate) Denegate() error {
	i, err := strconv.ParseInt(fmt.Sprintf("%v", p)[1:], 10, 0)
	if err != nil {
		return ilog.CheckErrAll(err, "error occured while denegating")()
	}
	p.assume(int(i))
	return nil
}

// MustDenegate turns a positive Coordinate integer into a negative one
func (p *Coordinate) MustDenegate() {
	if err := p.Denegate(); err != nil {
		panic(err)
	}
}

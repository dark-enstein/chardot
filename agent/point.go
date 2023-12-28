package agent

import (
	"context"
	"math"
)

// Point is stateful. The current position of Agent as A result of all the travels thus far
type Point struct {
	X, Y Coordinate
	ctx  context.Context
}

// distance calculates the distance between two points using Euclidean distance formula
func (p1 *Point) distance(p2 *Point) float64 {
	return math.Sqrt(math.Pow(float64(p2.X-p1.X), 2) + math.Pow(float64(p2.Y-p1.Y), 2))
}

func (d *Point) Path() *Path {
	var dist = NewPath(Dimensions)
	if d.X < 0 {
		pace := NewPace(LEFT)
		pace.ScalarMove(d.X)
		dist.A[0] = *pace
		dist.M[0] = *pace.PMap()
	} else if d.X > 0 {
		pace := NewPace(RIGHT)
		pace.ScalarMove(d.X)
		dist.A[0] = *pace
		dist.M[0] = *pace.PMap()
	}

	if d.Y < 0 {
		pace := NewPace(BACKWARD)
		pace.ScalarMove(d.Y)
		dist.A[1] = *pace
		dist.M[1] = *pace.PMap()
	} else if d.Y > 0 {
		pace := NewPace(FORWARD)
		pace.ScalarMove(d.Y)
		dist.A[1] = *pace
		dist.M[1] = *pace.PMap()
	}
	return dist
}

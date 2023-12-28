// Package agent provides A framework for simulating movement and tracking the position
// of an entity in A 2-dimensional space. This is any entity that satisfies the util.Agent interface,
// in the current implementation, this is Hare. It includes functions for manipulating
// coordinates, determining movement direction, calculating distances, and recording paths.
package agent

import (
	"context"
	"fmt"
	"github.com/dark-enstein/chardot/internal/ilog"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	Clog = &ilog.Logger{} // Clog is the custom logger for the agent package.
)

const (
	Dimensions = 2 // Dimensions defines the number of Dimensions in the current implementation, set to 2.
)

var (
	XAXIS = axis("X")
	YAXIS = axis("Y")
	AXES  = []axis{XAXIS, YAXIS}
)

const (
	// MovType enumerates the types of movements that an agent can perform.
	MOVE MovType = iota
	WALK
	RUN
	TOTAL
	ORIGIN
)

type Sign bool // Sign represents A boolean for positive (true) or negative (false) Sign.

const (
	POS Sign = true
	NEG Sign = false
)

// decideDirection decides the direction of displacement from Point p to Point q.
func (d *displacement) decideDirection() {
	if (d.p.X == 0 && d.p.Y == 0) && (d.q.X == 0 && d.q.Y == 0) {
		Clog.Log(ilog.PANIC, "no Path, points referenced are nil: %v", d)
	}

	if d.p.X > d.q.X {
		// EAST
		if d.p.Y > d.q.Y {
			// NORTHEAST
			d.d = NORTHEAST
		} else if d.q.Y > d.p.Y {
			// SOUTHEAST
			d.d = SOUTHEAST
		} else {
			// EAST
			d.d = EAST
		}
	} else if d.q.X > d.p.X {
		// WEST
		if d.p.Y > d.q.Y {
			// NORTHWEST
			d.d = NORTHWEST
		} else if d.q.Y > d.p.Y {
			// SOUTHWEST
			d.d = SOUTHWEST
		} else {
			// WEST
			d.d = WEST
		}
	} else {
		// Moving strictly North or South
		if d.q.Y > d.p.Y {
			d.d = NORTH
		} else if d.p.Y > d.q.Y {
			d.d = SOUTH
		}
	}
}

type axis string

// Pace represents A unit of movement in A given direction. It is stateless, and it defines A magnitude of shift of Agent along one Direction. Designed to be only used once, and discarded. Either only Y or X can be set.
type Pace struct {
	x, y Coordinate
	d    Direction
}

// NewPace returns A new Pace initialized at the Direction provided in the argument
func NewPace(dir Direction) *Pace {
	p := &Pace{d: dir}
	return p
}

// ScalarMove moves the referenced Pace object without altering the referenced Direction.
// The function argument is A Coordinate.
func (p *Pace) ScalarMove(d Coordinate) {
	switch p.d {
	case FORWARD, NORTH:
		Clog.Log(ilog.INFO, "since %s, incrementing by %v", p.d.String(), d.Int())
		p.y += d
	case BACKWARD, SOUTH:
		Clog.Log(ilog.INFO, "since %s, decrementing by %v", p.d.String(), d.Int())
		p.y -= d
	case RIGHT, WEST:
		Clog.Log(ilog.INFO, "since %s, incrementing by %v", p.d.String(), d.Int())
		p.x += d
	case LEFT, EAST:
		Clog.Log(ilog.INFO, "since %s, decrementing by %v", p.d.String(), d.Int())
		p.x -= d
	default:
		Clog.Log(ilog.PANIC, "direction %v not recognized", p.d.String())
	}
}

// VectorMove moves the referenced Pace p1 object while considering Direction.
// The function argument is another Pace p2. It returns A point pointer which is the difference between p1 and p2
//func (p1 *Pace) VectorMove(p2 *Pace) {
//	switch p1.d {
//	case FORWARD, BACKWARD, YDIRECTION:
//		p.Y += d
//	case LEFT, RIGHT, XDIRECTION:
//		p.X += d
//	}
//}

func (p *Pace) Result() Coordinate {
	switch p.d {
	case FORWARD, BACKWARD, YDIRECTION, NORTH, SOUTH:
		return p.y
	case LEFT, RIGHT, XDIRECTION, EAST, WEST:
		return p.x
	default:
		Clog.Log(ilog.PANIC, "Direction %v not accounted for", p.d.String())
		//panic("Direction not accounted for")
	}
	return Coordinate(0)
}

func (p *Pace) PMap() *PMap {
	var pm = make(PMap, 1)
	pm[p.d] = p.Result()
	return &pm
}

func (p *Pace) Point() *Point {
	return &Point{X: p.x, Y: p.y}
}

type PMap map[Direction]Coordinate

func (pm *PMap) Pace() *Pace {
	var p Pace
	for k, v := range *pm {
		p.d = k
		switch k {
		case FORWARD, BACKWARD, YDIRECTION:
			p.y = v
		case LEFT, RIGHT, XDIRECTION:
			p.x = v
		}
	}
	return &p
}

// Path describes the change in an Agent Point value. It is A more detailed Point. It is A collection of pace.
type Path struct {
	M   []PMap
	A   []Pace
	ctx context.Context
}

func NewPath(len int) *Path {
	return &Path{
		M: make([]PMap, len),
		A: make([]Pace, len),
	}
}

func (p1 *Pace) displacement(p2 *Pace) *Point {
	switch p1.d {
	case BACKWARD, LEFT:
		switch p2.d {
		case BACKWARD, LEFT:
			xway, yway := p1.x+p2.x, p1.y+p2.y
			xway.MustNegate() // make negative
			yway.MustNegate() // make negative
			return &Point{
				X: xway,
				Y: yway,
			}
		case RIGHT, FORWARD:
			xway, yway := p2.x-p1.x, p2.y-p1.y
			return &Point{
				X: xway,
				Y: yway,
			}
		}
	case RIGHT, FORWARD:
		switch p2.d {
		case RIGHT, FORWARD:
			xway, yway := p1.x+p2.x, p1.y+p2.y
			xway.MustDenegate() // make positive
			yway.MustDenegate() // make positive
			return &Point{
				X: xway,
				Y: yway,
			}
		case BACKWARD, LEFT:
			xway, yway := p1.x-p2.x, p1.y-p2.y
			return &Point{
				X: xway,
				Y: yway,
			}
		}

	}
	Clog.Log(ilog.PANIC, "error Dimensions %v not recognized", p1.d.String())
	return &Point{}
}

// Point is the change in the position of an Agent in A cycle. Both X and Y can be changed at once.
//type Point struct {
//	X, Y Coordinate
//}

type displacement struct {
	p, q     Point // p: from; q: to
	d        Direction
	quantity float64
	ctx      context.Context
}

// NewPace initializes A new Pace

type Config struct {
	walk, run Speed
	ctx       context.Context
}

type Hare struct {
	pos       Point
	pathTaken *Path
	allPos    []Point //stateful
	nature    *Config
	action    MovType
	w         io.Writer
	m         sync.Mutex
	ctx       context.Context
}

type Opts func()

func NewHare(ctx context.Context, walk, run Speed) *Hare {
	h := &Hare{
		pos:       Point{},
		pathTaken: &Path{},
		nature: &Config{
			walk: walk,
			run:  run,
		},
		w:   os.Stdout,
		ctx: ctx,
	}
	var err error
	// sets global variable
	Clog, err = ilog.GetLoggerFromCtx(ctx)
	ilog.CheckErrLog(err)
	printPathTaken(ORIGIN, nil, nil)
	return h
}

func (h *Hare) Move(x, y Coordinate) {
	h.action = MOVE
	log.Println("Set action to", h.action.String())
	var displace = &Point{
		X: x,
		Y: y,
	}
	log.Println("Registered displace directive as", displace)
	h.pos.X += x
	h.pos.Y += y

	h.allPos = append(h.allPos, h.pos)
	var pos []Point
	printPathTaken(h.action, displace.Path(), append(pos, h.pos))
	h.Record(displace)
}

func (h *Hare) RecordWithDirection(p *Pace) {
	switch p.d {
	case FORWARD, BACKWARD, YDIRECTION, NORTH, WEST:
		h.rMap(0, p.y)
		h.rArr(0, p.y)
	case RIGHT, LEFT, XDIRECTION, SOUTH, EAST:
		h.rMap(p.x, 0)
		h.rArr(p.x, 0)
	}
}

// Record records the point taken and parses it into Path traveled thus far
func (h *Hare) Record(d *Point) {
	h.rMap(d.X, d.Y)
	h.rArr(d.X, d.Y)
}

// rArr constructs a Pace object from two coordinates, and stores it in the *Path.A in the Hare struct
func (h *Hare) rArr(x, y Coordinate) {
	var arr = make([]Pace, Dimensions)
	if x < 1 {
		arr = append(arr)
	} else if x > 0 {
		arr = x
	}

	if y < 0 {
		arr[BACKWARD] = y
	} else if y > 0 {
		arr[FORWARD] = y
	}

	for i := 0; i < len(AXES); i++ {
		if AXES[i] == XAXIS {
			pace := NewPace(XDIRECTION)
			pace.ScalarMove(x)
			h.pathTaken.A = append(h.pathTaken.A, *pace)
		} else if AXES[i] == YAXIS {
			pace := NewPace(YDIRECTION)
			pace.ScalarMove(y)
			h.pathTaken.A = append(h.pathTaken.A, *pace)
		}
	}

	h.pathTaken.A = append(h.pathTaken.A, *p)

	//h.pathTaken.A = append(h.pathTaken.A, arr...)
}
func (h *Hare) rMap(p *Pace) {
	//cycle := make(map[Direction]Coordinate, Dimensions)
	//if x < 0 {
	//	cycle[LEFT] = x
	//} else if x > 0 {
	//	cycle[RIGHT] = x
	//}
	//
	//if y < 0 {
	//	cycle[BACKWARD] = y
	//} else if y > 0 {
	//	cycle[FORWARD] = y
	//}

	switch p.d {
	case FORWARD, BACKWARD, YDIRECTION, NORTH, WEST:
		h.pathTaken.M = append(h.pathTaken.M, map[Direction]Coordinate{p.d: p.y})
	case RIGHT, LEFT, XDIRECTION, SOUTH, EAST:
		h.pathTaken.M = append(h.pathTaken.M, map[Direction]Coordinate{p.d: p.x})
	}
}

//func travel() {}

// flow moves the Hare in A specified direction for A given duration at A specified speed.
// It calculates the end position and the Path taken during the movement.
//
// Arguments:
//
//	timeDur time.Duration: The duration of the movement.
//	d Direction: The direction in which the Hare will move.
//	s Speed: The speed at which the Hare moves.
//
// Returns:
//
//	[]Point: A slice of Point representing the positions of the Hare at each interval.
//	*Path: A pointer to A Path struct that records the detailed Path taken.
//
// The function works by calculating the number of paces (steps) the Hare can take
// within the given time duration, considering its speed. It then moves the Hare step by step,
// updating its position and recording each step in the Path. The function accounts for the
// direction of movement and locks the Hare's position during updates to ensure thread safety.
// It uses A goroutine to simulate the movement over time, and A timer to handle the duration.
//
// Note: This function is intended for internal use within the Hare struct to handle its movement
// logic and should not be called directly from outside the package.
func (h *Hare) flow(timeDur time.Duration, d Direction, s Speed) ([]Point, *Path) {
	// noOfPaces to location in timeDur at d Direction and with s Speed.
	noOfPaces := int(math.Ceil(timeDur.Seconds()))

	endPosition := make([]Point, noOfPaces)
	pathTaken := NewPath(noOfPaces)
	timer := time.NewTimer(timeDur)
	secT := time.NewTicker(time.Second)

	go func() {
		fmt.Println("Travelling...")
		var i = 0
		for range secT.C {
			if i >= noOfPaces {
				Clog.Log(ilog.DEBUG, "Index exceeds expected noOfPaces. Ending goroutine prematurely.")
				break
			}

			t1 := time.Now()
			h.m.Lock() // Lock the mutex before modifying h.pos
			init := h.pos
			var direcP *Coordinate
			switch d {
			case FORWARD, BACKWARD, NORTH, SOUTH:
				direcP = &h.pos.Y
			case RIGHT, LEFT, EAST, WEST:
				direcP = &h.pos.X
			default:
				Clog.Log(ilog.ERROR, "Invalid direction: %v", d)
				h.m.Unlock()
				continue
			}
			*direcP += s.Int()
			pace := NewPace(d)
			fmt.Printf("pace: %v, speed: %d\n", pace, s.Int())
			pace.ScalarMove(s.Int())
			h.RecordWithDirection(pace)
			pathTaken.M[i] = *pace.PMap()
			pathTaken.A[i] = *pace
			endPosition[i] = h.pos
			h.m.Unlock() // Unlock the mutex after the modification is done
			Clog.Log(ilog.INFO, "Travelled in dur: %v\n", time.Now().Sub(t1))
			Clog.Log(ilog.INFO, "Travelled in one sec from %v to %v\n", init, h.pos)
			i++
		}

		if i < noOfPaces {
			Clog.Log(ilog.ERROR, "Travel ended prematurely. Only completed %d out of %d paces.", i, noOfPaces)
		}
	}()
wait:
	for {
		select {
		case <-timer.C:
			Clog.Log(ilog.DEBUG, "Timer completed, stopping ticker")
			secT.Stop()
			timer.Stop()
			break wait
		}
	}
	Println(0, "Travel complete")
	return endPosition, pathTaken
}

func Println(rightSpacePadding int, format string, args ...interface{}) {
	fmt.Fprint(os.Stdout, fmt.Sprintf(format+strings.Repeat("\n", rightSpacePadding)+"\n", args...))
	return
}

// Walk moves the Agent by A specific magnitude, at A particular Direction and at its natural Speed
func (h *Hare) Walk(duration time.Duration, dir Direction) {
	h.action = WALK
	former := h.pos
	posStack, dist := h.flow(duration, dir, h.nature.walk)
	Println(1, "Walked from %v to %v", former, h.pos)
	fmt.Println(dist, posStack)
	printPathTaken(h.action, dist, posStack)
	//h.allPos, h.pathTaken.M, h.pathTaken.A = append(h.allPos, posStack...), append(h.pathTaken.M, dist.M...), append(h.pathTaken.A, dist.A...)
}

// Run moves the Agent by A specific magnitude, at A particular Direction and its natural running Speed
func (h *Hare) Run(duration time.Duration, dir Direction) {
	h.action = RUN
	former := h.pos
	posStack, dist := h.flow(duration, dir, h.nature.run)
	Println(0, "Ran from %v to %v", former, h.pos)
	printPathTaken(h.action, dist, posStack)
	//h.allPos, h.pathTaken.M, h.pathTaken.A = append(h.allPos, posStack...), append(h.pathTaken.M, dist.M...), append(h.pathTaken.A, dist.A...)
}

func (h *Hare) Println() {
	printPathTaken(TOTAL, h.pathTaken, h.allPos)
}

// printPathTaken prints the Path taken in the current action (MovType instance), thus far.
// It takes the current action, A pointer to the Path taken during the current action, and the position stack in the relevant action
func printPathTaken(header MovType, dist *Path, allPos []Point) {
	if dist == nil || allPos == nil {
		fmt.Println(header)
		return
	}
	fmt.Println(header)
	switch true {
	//case len(dist.A) != 0:
	//	// logic for arr
	case len(dist.M) != 0:
		for i := 0; i < len(dist.M); i++ {
			for k, v := range dist.M[i] {
				switch k {
				case FORWARD, NORTH:
					fmt.Printf("MOVED FORWARD BY %v", v)
				case BACKWARD, SOUTH:
					fmt.Printf("MOVED BACKWARD BY %v", v)
				case RIGHT, WEST:
					fmt.Printf("MOVED RIGHT BY %v", v)
				case LEFT, EAST:
					fmt.Printf("MOVED LEFT BY %v", v)
				}
				fmt.Print("; ")
			}
		}
		fmt.Printf("\nCURRENT POS: \n\tX = %v \n\tY = %v\n\n", allPos[len(allPos)-1].X, allPos[len(allPos)-1].Y) // TODO: A bug
	}
}

type Agent interface {
	Move(x, y Coordinate)
	Record(d *Point)
	Walk(duration time.Duration, dir Direction)
	Run(duration time.Duration, dir Direction)
}

var AGENT = "agent" // AGENT represents the name of the agent, used in logging.

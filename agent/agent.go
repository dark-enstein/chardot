package agent

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	dimensions = 2
)

var (
	XAXIS = axis("x")
	YAXIS = axis("y")
	AXES  = []axis{XAXIS, YAXIS}
)

const (
	MOVE movType = iota
	WALK
	RUN
	TOTAL
	ORIGIN
)

type sign bool

const (
	POS sign = true
	NEG sign = false
)

// coordinate is a wrapper over the standard int. it helps for adding convenient functionalities such as negate()
type coordinate int

// Int returns the underlying integer value of a coordinate
func (p *coordinate) Int() int {
	return int(*p)
}

// assume overrides the integer value of a coordinate
func (p *coordinate) assume(i int) {
	point := coordinate(i)
	p = &point
}

func (p *coordinate) add(i int) {
	po := coordinate(i)
	*p += po
}

// Sign extracts the underlying sign from a coordinate
func (p *coordinate) Sign() sign {
	if strconv.FormatInt(int64(*p), 10)[0] != '-' {
		return POS
	}
	return NEG
}

// negate turns a positive coordinate integer into a negative one
func (p *coordinate) negate() error {
	i, err := strconv.ParseInt(fmt.Sprintf("-%v", p), 10, 0)
	if err != nil {
		log.Println("error occured while negating")
		return err
	}
	p.assume(int(i))
	return nil
}

// mustNegate turns a positive coordinate integer into a negative one
func (p *coordinate) mustNegate() {
	if err := p.negate(); err != nil {
		panic(err)
	}
}

// denegate turns a negative coordinate integer into a positive one
func (p *coordinate) denegate() error {
	i, err := strconv.ParseInt(fmt.Sprintf("%v", p)[1:], 10, 0)
	if err != nil {
		log.Println("error occured while denegating")
		return err
	}
	p.assume(int(i))
	return nil
}

// mustDenegate turns a positive coordinate integer into a negative one
func (p *coordinate) mustDenegate() {
	if err := p.denegate(); err != nil {
		panic(err)
	}
}

// Point estimates the difference between two points
func (p1 *Point) displacement(p2 *Point) *displacement {
	var displace displacement
	displace.p, displace.q = *p1, *p2
	displace.quantity = p1.distance(p2)
	displace.decideDirection()
	return &displace
}

// decideDirection decides the direction of p -> q based on their points coordinates
func (d *displacement) decideDirection() {
	if (d.p.x == 0 && d.p.y == 0) && (d.q.x == 0 && d.q.y == 0) {
		log.Panicln("no Path, points referenced are nil")
	}

	if d.p.x > d.q.x {
		// EAST
		if d.p.y > d.q.y {
			// NORTHEAST
			d.d = NORTHEAST
		} else if d.q.y > d.p.y {
			// SOUTHEAST
			d.d = SOUTHEAST
		} else {
			// EAST
			d.d = EAST
		}
	}

	if d.q.x > d.p.x {
		// WEST
		if d.p.y > d.q.y {
			// NORTHWEST
			d.d = NORTHWEST
		} else if d.q.y > d.p.y {
			// SOUTHWEST
			d.d = SOUTHWEST
		} else {
			// WEST
			d.d = WEST
		}
	}
}

// distance calculates the distance between two points using Euclidean distance formula
func (p1 *Point) distance(p2 *Point) float64 {
	return math.Sqrt(math.Pow(float64(p2.x-p1.x), 2) + math.Pow(float64(p2.y-p1.y), 2))
}

// movType defines the type of action being carried out on an Agent
type movType int

// String extracts the type of the referencing action as a string
func (m movType) String() string {
	switch m {
	case MOVE:
		return fmt.Sprintf("MOVE")
	case WALK:
		return fmt.Sprintf("WALK")
	case RUN:
		return fmt.Sprintf("RUN")
	case TOTAL:
		return fmt.Sprintf("TOTAL")
	case ORIGIN:
		return fmt.Sprintf("ORIGIN\n")
	}
	return "action unrecognized"
}

type axis string

type direction int

const (
	FORWARD direction = iota
	BACKWARD
	YDIRECTION
	LEFT
	RIGHT
	XDIRECTION
	NORTHEAST
	NORTHWEST
	SOUTHEAST
	SOUTHWEST
	NORTH
	SOUTH
	EAST
	WEST
)

// Pace is stateless, and it defines a magnitude of shift of Agent along one direction. Designed to be only used once, and discarded. Either only y or x can be set.
type Pace struct {
	x, y coordinate
	d    direction
}

// NewPace returns a new Pace initialized at the direction provided in the argument
func NewPace(dir direction) *Pace {
	p := &Pace{d: dir}
	return p
}

// ScalarMove moves the referenced Pace object without altering the referenced direction.
// The function argument is a coordinate.
func (p *Pace) ScalarMove(d coordinate) {
	switch p.d {
	case FORWARD, BACKWARD, YDIRECTION:
		p.y += d
	case LEFT, RIGHT, XDIRECTION:
		p.x += d
	}
}

// VectorMove moves the referenced Pace p1 object while considering direction.
// The function argument is another Pace p2. It returns a point pointer which is the difference between p1 and p2
//func (p1 *Pace) VectorMove(p2 *Pace) {
//	switch p1.d {
//	case FORWARD, BACKWARD, YDIRECTION:
//		p.y += d
//	case LEFT, RIGHT, XDIRECTION:
//		p.x += d
//	}
//}

func (p *Pace) Result() coordinate {
	switch p.d {
	case FORWARD, BACKWARD, YDIRECTION:
		return p.y
	case LEFT, RIGHT, XDIRECTION:
		return p.x
	default:
		log.Panicf("direction %v not accounted for", p.d)
		panic("direction not accounted for")
	}
}

func (p *Pace) PMap() *PMap {
	var pm = make(PMap, 1)
	pm[p.d] = p.Result()
	return &pm
}

func (p *Pace) Point() *Point {
	return &Point{x: p.x, y: p.y}
}

type PMap map[direction]coordinate

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

// path describes the change in an Agent Point value. It is a more detailed Point. It is a collection of pace.
type path struct {
	m []PMap
	a []Pace
}

func newPath(len int) *path {
	return &path{
		m: make([]PMap, len),
		a: make([]Pace, len),
	}
}

func (p1 *Pace) displacement(p2 *Pace) *Point {
	switch p1.d {
	case BACKWARD, LEFT:
		switch p2.d {
		case BACKWARD, LEFT:
			xway, yway := p1.x+p2.x, p1.y+p2.y
			xway.mustNegate() // make negative
			yway.mustNegate() // make negative
			return &Point{
				x: xway,
				y: yway,
			}
		case RIGHT, FORWARD:
			xway, yway := p2.x-p1.x, p2.y-p1.y
			return &Point{
				x: xway,
				y: yway,
			}
		}
	case RIGHT, FORWARD:
		switch p2.d {
		case RIGHT, FORWARD:
			xway, yway := p1.x+p2.x, p1.y+p2.y
			xway.mustDenegate() // make positive
			yway.mustDenegate() // make positive
			return &Point{
				x: xway,
				y: yway,
			}
		case BACKWARD, LEFT:
			xway, yway := p1.x-p2.x, p1.y-p2.y
			return &Point{
				x: xway,
				y: yway,
			}
		}

	}
	log.Panicln("error dimensions not recognized")
	return &Point{}
}

// Point is the change in the position of an Agent in a cycle. Both X and Y can be changed at once.
//type Point struct {
//	x, y coordinate
//}

type displacement struct {
	p, q     Point // p: from; q: to
	d        direction
	quantity float64
}

func (d *Point) Path() *path {
	var dist = newPath(dimensions)
	if d.x < 0 {
		pace := NewPace(LEFT)
		pace.ScalarMove(d.x)
		dist.a[0] = *pace
		dist.m[0] = *pace.PMap()
	} else if d.x > 0 {
		pace := NewPace(RIGHT)
		pace.ScalarMove(d.x)
		dist.a[0] = *pace
		dist.m[0] = *pace.PMap()
	}

	if d.y < 0 {
		pace := NewPace(BACKWARD)
		pace.ScalarMove(d.y)
		dist.a[1] = *pace
		dist.m[1] = *pace.PMap()
	} else if d.y > 0 {
		pace := NewPace(FORWARD)
		pace.ScalarMove(d.y)
		dist.a[1] = *pace
		dist.m[1] = *pace.PMap()
	}
	return dist
}

// NewPace initializes a new Pace

type Agent interface {
	Move(x, y coordinate)
	Record(d *Point)
}

// Point is stateful. The current position of Agent as a result of all the travels thus far
type Point struct {
	x, y coordinate
}

func (p *Point) MoveBy(q *Point) *Point {
	if p.x > q.x {

	}
	return &Point{}
}

type speed int

func (s speed) Int() coordinate {
	return coordinate(s)
}

type Config struct {
	walk, run speed
}

type Hare struct {
	pos       Point
	pathTaken *path
	allPos    []Point //stateful
	nature    *Config
	action    movType
	w         io.Writer
	m         sync.Mutex
}

func NewHare(walk, run speed) *Hare {
	h := &Hare{
		pos:       Point{},
		pathTaken: &path{},
		nature: &Config{
			walk: walk,
			run:  run,
		},
		w: os.Stdout,
	}
	printPathTaken(ORIGIN, nil, nil)
	return h
}

func (h *Hare) Move(x, y coordinate) {
	h.action = MOVE
	log.Println("Set action to", h.action.String())
	var displace = &Point{
		x: x,
		y: y,
	}
	log.Println("Registered displace directive as", displace)
	h.pos.x += x
	h.pos.y += y

	h.allPos = append(h.allPos, h.pos)
	var pos []Point
	printPathTaken(h.action, displace.Path(), append(pos, h.pos))
	h.Record(displace)
}

func (h *Hare) RecordWithDirection(p *Pace) {
	switch p.d {
	case FORWARD, BACKWARD, YDIRECTION:
		h.rMap(0, p.y)
		h.rArr(0, p.y)
	case RIGHT, LEFT, XDIRECTION:
		h.rMap(p.x, 0)
		h.rArr(p.x, 0)
	}
}

// Record records the point taken and parses it into path traveled thus far
func (h *Hare) Record(d *Point) {
	h.rMap(d.x, d.y)
	h.rArr(d.x, d.x)
}

func (h *Hare) rArr(x, y coordinate) {
	var arr = make([]Pace, dimensions)
	for i := 0; i < len(AXES); i++ {
		if AXES[i] == XAXIS {
			pace := NewPace(XDIRECTION)
			pace.ScalarMove(x)
			arr = append(arr, *pace)
		} else if AXES[i] == YAXIS {
			pace := NewPace(YDIRECTION)
			pace.ScalarMove(y)
			arr = append(arr, *pace)
		}
	}

	h.pathTaken.a = append(h.pathTaken.a, arr...)
}
func (h *Hare) rMap(x, y coordinate) {
	cycle := make(map[direction]coordinate, dimensions)
	if x < 0 {
		cycle[LEFT] = x
	} else if x > 0 {
		cycle[RIGHT] = x
	}

	if y < 0 {
		cycle[BACKWARD] = y
	} else if y > 0 {
		cycle[FORWARD] = y
	}

	h.pathTaken.m = append(h.pathTaken.m, cycle)
}

//func travel() {}

// flow moves the Point in only one direction.
func (h *Hare) flow(timeDur time.Duration, d direction, s speed) ([]Point, *path) {
	// noOfPaces to location in timeDur at d direction and with s speed.
	noOfPaces := int(math.Ceil(timeDur.Seconds()))

	endPosition := make([]Point, noOfPaces)
	pathTaken := newPath(noOfPaces)

	var direcP *coordinate
	switch d {
	case FORWARD, BACKWARD:
		direcP = &h.pos.y
	case RIGHT, LEFT:
		direcP = &h.pos.x
	}

	timer := time.NewTimer(timeDur)

	secT := time.NewTicker(time.Second)
	go func() {
		fmt.Println("Travelling...")
		var i = 0
		for _ = range secT.C {
			t1 := time.Now()
			init := h.pos
			*direcP += s.Int()
			pace := NewPace(d)
			pace.ScalarMove(s.Int())
			h.m.Lock()
			h.RecordWithDirection(pace)
			pathTaken.m[i] = *pace.PMap()
			pathTaken.a[i] = *pace
			endPosition[i] = h.pos
			h.m.Unlock()
			log.Printf("Travelled in dur: %v\n", time.Now().Sub(t1))
			log.Printf("Travelled in one sec from %v to %v\n", init, h.pos)
			i++
		}
	}()
wait:
	for {
		select {
		case <-timer.C:
			log.Println("Timer completed, stopping ticker")
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

// Walk moves the Agent by a specific magnitude, at a particular direction and at its natural speed
func (h *Hare) Walk(duration time.Duration, dir direction) {
	h.action = WALK
	former := h.pos
	posStack, dist := h.flow(duration, dir, h.nature.walk)
	Println(1, "Walked from %v to %v", former, h.pos)
	fmt.Println(dist, posStack)
	printPathTaken(h.action, dist, posStack)
	//h.allPos, h.pathTaken.m, h.pathTaken.a = append(h.allPos, posStack...), append(h.pathTaken.m, dist.m...), append(h.pathTaken.a, dist.a...)
}

// Run moves the Agent by a specific magnitude, at a particular direction and its natural running speed
func (h *Hare) Run(duration time.Duration, dir direction) {
	h.action = RUN
	former := h.pos
	posStack, dist := h.flow(duration, dir, h.nature.run)
	Println(0, "Ran from %v to %v", former, h.pos)
	printPathTaken(h.action, dist, posStack)
	h.allPos, h.pathTaken.m, h.pathTaken.a = append(h.allPos, posStack...), append(h.pathTaken.m, dist.m...), append(h.pathTaken.a, dist.a...)
}

func (h *Hare) Println() {
	printPathTaken(TOTAL, h.pathTaken, h.allPos)
}

// printPathTaken prints the path taken in the current action (movType instance), thus far.
// It takes the current action, a pointer to the path taken during the current action, and the position stack in the relevant action
func printPathTaken(header movType, dist *path, allPos []Point) {
	if dist == nil || allPos == nil {
		fmt.Println(header)
		return
	}
	fmt.Println(header)
	switch true {
	//case len(dist.a) != 0:
	//	// logic for arr
	case len(dist.m) != 0:
		for i := 0; i < len(dist.m); i++ {
			for k, v := range dist.m[i] {
				switch k {
				case FORWARD:
					fmt.Printf("MOVED FORWARD BY %v", v)
				case BACKWARD:
					fmt.Printf("MOVED BACKWARD BY %v", v)
				case RIGHT:
					fmt.Printf("MOVED RIGHT BY %v", v)
				case LEFT:
					fmt.Printf("MOVED LEFT BY %v", v)
				}
				fmt.Print("; ")
			}
		}
		fmt.Printf("\nCURRENT POS: \n\tX = %v \n\tY = %v\n\n", allPos[len(allPos)-1].x, allPos[len(allPos)-1].y) // TODO: a bug
	}
}

package agent

import (
	"fmt"
	"github.com/dark-enstein/chardot/internal/ilog"
)

// MovType defines the type of action being carried out on an Agent
type MovType int

// String extracts the type of the referencing action as a string
func (m MovType) String() string {
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

type Direction int

const (
	FORWARD Direction = iota
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

func (d Direction) String() string {
	switch d {
	case FORWARD:
		return "FORWARD"
	case BACKWARD:
		return "BACKWARD"
	case YDIRECTION:
		return "YDIRECTION"
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	case XDIRECTION:
		return "XDIRECTION"
	case NORTHEAST:
		return "NORTHEAST"
	case NORTHWEST:
		return "NORTHWEST"
	case SOUTHEAST:
		return "SOUTHEAST"
	case SOUTHWEST:
		return "SOUTHWEST"
	case NORTH:
		return "NORTH"
	case SOUTH:
		return "SOUTH"
	case EAST:
		return "EAST"
	case WEST:
		return "WEST"
	default:
		Clog.Log(ilog.PANIC, "direction %v unrecognized", d)
	}
	return ""
}

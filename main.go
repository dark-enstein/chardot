package main

import (
	"github.com/dark-enstein/chardot/agent"
	"time"
)

// Track user location across 2d space
// Encode functions to run, walk, and wait

func main() {
	h := agent.NewHare(4, 6)

	h.Move(4, 5)

	h.Move(10, -2)

	h.Walk(time.Second*6, agent.RIGHT)
	h.Run(time.Second*10, agent.LEFT)

	h.Println()
}

package main

import (
	"github.com/dark-enstein/chardot/agent"
	"github.com/dark-enstein/chardot/cfg"
	"github.com/dark-enstein/chardot/internal/ilog"
	"time"
)

// Track user location across 2d space
// Encode functions to run, walk, and wait

func main() {
	c := cfg.Config{
		WalkSpeed: "4",
		RunSpeed:  "6",
	}
	ctx, _ := c.InitSetUp()
	_, err := ilog.GetLoggerFromCtx(ctx)
	ilog.CheckErrLog(err)
	//clog.Log(ilog.PANIC, "errors encountered during init: %v", errs)
	h := agent.NewHare(ctx, 6, 6)

	h.Move(4, 5)

	h.Move(10, -2)

	h.Walk(time.Second*6, agent.RIGHT)
	h.Run(time.Second*10, agent.LEFT)

	h.Println()
}

package cfg

import (
	"context"
	"fmt"
	"github.com/dark-enstein/chardot/agent"
	"github.com/dark-enstein/chardot/internal/ilog"
	"github.com/dark-enstein/chardot/util"
	"log"
	"time"
)

var (
	ERRSPEEDNOTDEFINED     = fmt.Errorf("Speed not defined for walk\n\n")
	ERRORNOTVALIDRETURNING = fmt.Errorf("LogLevel passed in invalid. Using INFO.")
	DEFAULTWALKSPEED       = agent.Speed(0)
	DEFAULTRUNSPEED        = agent.Speed(0)
)

type Config struct {
	A         []Action `yaml:"actions"`
	LogLevel  string   `yaml:"logLevel"`
	WalkSpeed string   `yaml:"walkSpeed"`
	RunSpeed  string   `yaml:"runSpeed"`
}

func NewConfig(loglevel, walkS, runS string, acts ...Action) *Config {
	log.Println("set up new config")
	return &Config{
		A:         acts,
		LogLevel:  loglevel,
		WalkSpeed: walkS,
		RunSpeed:  runS,
	}
}

type Configurer interface {
	SetUp() func() []error
}

func (c *Config) InitSetUp() (context.Context, error) {
	log.Println("initializing setup")
	if c.LogLevel == "" {
		c.LogLevel = "INFO"
	}

	if c.LogLevel != "INFO" && c.LogLevel != "DEBUG" && c.LogLevel != "ERROR" && c.LogLevel != "PANIC" {
		log.Println(ERRORNOTVALIDRETURNING)
		c.LogLevel = "INFO"
	}

	logger, err := ilog.NewLogger(c.LogLevel)
	if err != nil {
		log.Println(err)
	}

	ctx := context.WithValue(context.Background(), ilog.LOGGERCTX, logger)

	//Values passed in:

	ag, err := c.SetUpAgent(ctx)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, agent.AGENT, ag)

	return ctx, err
}

func (c *Config) SetUp() error {
	ctx, err := c.InitSetUp()

	if err != nil {
		return err
	}
	var ext []Command
	for i := 0; i < len(c.A); i++ {
		cmd, err := c.A[i].IntoCommand()
		if err != nil {
			break
		}
		ext = append(ext, cmd)
	}
	if err != nil {
		return err
	}
	return func() error {
		for i := 0; i < len(ext); i++ {
			err := ext[i].Do(ctx)
			if err != nil {
				break
			}
		}
		return err
	}()
}

func (c *Config) ResolveSpeed(ctx context.Context) (w, r *agent.Speed, err error) {
	clog, err := ilog.GetLoggerFromCtx(ctx)
	ilog.CheckErrLog(err)
	if c.WalkSpeed == "" {
		clog.Log(ilog.DEBUG, "WalkSpeed is not defined in the configuration. Defaulting to %v", DEFAULTWALKSPEED)
	}

	if c.RunSpeed == "" {
		clog.Log(ilog.DEBUG, "RunSpeed is not defined in the configuration. Defaulting to %v", DEFAULTRUNSPEED)
	}
	return agent.Speed(util.MustAtoi(c.WalkSpeed)).Ptr(), agent.Speed(util.MustAtoi(c.RunSpeed)).Ptr(), err
}

func (c *Config) SetUpAgent(ctx context.Context) (agent.Agent, error) {
	var err error = nil
	walkS, runS, err := c.ResolveSpeed(ctx)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, agent.WALK, walkS)
	ctx = context.WithValue(ctx, agent.RUN, runS)
	walkS, err = agent.GetSpeedFromCtx(ctx, agent.WALK)
	if err != nil {
		return nil, err
	}
	runS, err = agent.GetSpeedFromCtx(ctx, agent.RUN)
	if err != nil {
		return nil, err
	}
	return agent.NewHare(ctx, *walkS, *runS), err
}

type Command interface {
	Do(ctx context.Context) error
}

type Action struct {
	Name        string `yaml:"name"`
	DurationSec int    `yaml:"duration"`
	Direction   string `yaml:"direction"`
}

func (a *Action) IntoCommand() (Command, error) {
	if a.DurationSec < 0 {
		return nil, fmt.Errorf("duration %v is negative. invalid", a.DurationSec)
	}
	direction := agent.Direction(-1)
	switch a.Direction {
	case "N":
		direction = agent.NORTH
	case "S":
		direction = agent.SOUTH
	case "E":
		direction = agent.EAST
	case "W":
		direction = agent.WEST
	default:
		return nil, fmt.Errorf("direction %v not recognized\n\n", a.Direction)
	}

	switch a.Name {
	case "walk":
		return &Walk{
			time:      time.Second * time.Duration(a.DurationSec),
			direction: direction,
			ctx:       nil,
		}, nil
	case "run":
		return &Run{
			time:      time.Second * time.Duration(a.DurationSec),
			direction: direction,
			ctx:       nil,
		}, nil
	}
	return nil, fmt.Errorf("action %v not recognized\n\n", a.Name)
}

type Walk struct {
	time      time.Duration
	direction agent.Direction
	ctx       context.Context
}

func (w *Walk) Do(ctx context.Context) error {
	ag, err := agent.GetAgentFromCtx(ctx)
	ilog.CheckErrLog(err)
	ag.Walk(w.time, w.direction)
	return nil
}

type Run struct {
	time      time.Duration
	direction agent.Direction
	ctx       context.Context
}

func (w *Run) Do(ctx context.Context) error {
	ag, err := agent.GetAgentFromCtx(ctx)
	ilog.CheckErrLog(err)
	ag.Run(w.time, w.direction)
	return nil
}

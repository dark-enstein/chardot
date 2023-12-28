package main

import (
	"fmt"
	"github.com/dark-enstein/chardot/cfg"
	"github.com/dark-enstein/chardot/internal/streams"
	"log"
	"os"
	"strings"
)

var (
	ERR_WRONGARG               = "err: wrong argument passed\n\n"
	ERR_VALNOTDEF              = "err: no value passed for arg\n\n"
	ERR_FILE404                = "err: config file referenced not found \n\n"
	ERR_CONFIGEXTINVALID       = "err: config file with invalid extension \n\n"
	ERR_CONFIGNOTINRIGHTFORMAT = "err: config file yaml isn't in the right format\n\n"
	ERR_TOOMANYFLAGS           = "err: too many flags passed in"
)

const (
	HELP = `
usage: chardot --file <file location>
Note: if --file isn't passed the default location .chardot.cfg is used

Cannot use this tool? Help us improve by raising an issue here: https://github.com/dark-enstein/chardot/issues/new
	`
	CONFIGEXT = "cfg"
)

var (
	ACCEPTEDFLGS = []string{"log_level", "file"}
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Println("you passed in no flags, running default config")
		errs := _runningDry()
		fmt.Println(errs)
	}
	switch len(args) {
	case 1:
		if args[0] != "--file" {
			_ = fmt.Errorf(ERR_WRONGARG)
			fmt.Println(HELP)
			os.Exit(1)
		}
		_ = fmt.Errorf(ERR_VALNOTDEF)
		os.Exit(1)
	case 2:
		if args[0] != "--file" {
			_ = fmt.Errorf(ERR_WRONGARG)
			fmt.Println(HELP)
			os.Exit(1)
		}
		fd, err := os.Stat(args[1])
		if err == os.ErrNotExist {
			fmt.Println(fmt.Errorf(ERR_FILE404))
			fmt.Println(HELP)
			os.Exit(1)
		}
		if err != nil {
			_ = fmt.Errorf(err.Error())
			fmt.Println(HELP)
			os.Exit(1)
		}

		slis := strings.Split(fd.Name(), ".")
		if slis[len(slis)-1] != CONFIGEXT {
			_ = fmt.Errorf(ERR_CONFIGEXTINVALID)
			fmt.Println(HELP)
			os.Exit(1)
		}

		data, err := os.ReadFile(args[1])
		if err != nil {
			_ = fmt.Errorf(err.Error())
			fmt.Println(HELP)
			os.Exit(1)
		}

		c, err := streams.YamlDecode(data)
		if err != nil {
			_ = fmt.Errorf(ERR_CONFIGNOTINRIGHTFORMAT)
			fmt.Println(HELP)
			os.Exit(1)
		}

		//run := c.SetUp()
		fmt.Println(c.SetUp())
	}

	// include flag for log_level
	// polish flags admission process, so it is easily extendable
	if len(args) > 2 {
		_ = fmt.Errorf(ERR_TOOMANYFLAGS)
		fmt.Println(HELP)
		os.Exit(1)
	}
	return
}

func _runningDry() error {
	w := []cfg.Action{
		{
			Name:        "walk",
			DurationSec: 5,
			Direction:   "N",
		},
		{
			Name:        "run",
			DurationSec: 50,
			Direction:   "E",
		},
	}
	c := cfg.NewConfig("INFO", "5", "6", w...)
	//run := c.SetUp()
	return c.SetUp()
}

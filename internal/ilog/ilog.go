package ilog

import (
	"context"
	"fmt"
	"log"
)

const (
	INFO int = iota
	DEBUG
	ERROR
	PANIC
)

var (
	ERRVARNOTRECOGNIZED = "%s not recognized\n\n"
	VAR_LOGLEVEL        = "log_level"
)

var (
	LOGGERCTX = "LOGGERCTX"
)

type Logger struct {
	level int
}

func NewLogger(level string) (*Logger, error) {
	switch level {
	case "INFO":
		return &Logger{level: INFO}, nil
	case "ERROR":
		return &Logger{level: ERROR}, nil
	case "DEBUG":
		return &Logger{level: DEBUG}, nil
	case "PANIC":
		return &Logger{level: PANIC}, nil
	default:
		return nil, fmt.Errorf(ERRVARNOTRECOGNIZED, VAR_LOGLEVEL)
	}
}

func (l *Logger) Log(lev int, msg string, args ...interface{}) {
	suffix := ""
	switch lev {
	case INFO:
		suffix = "info: "
	case ERROR:
		suffix = "error: "
	case DEBUG:
		suffix = "debug: "
	case PANIC:
		suffix = "panic: "
	}
	if lev >= l.level {
		switch lev {
		case INFO, ERROR, DEBUG:
			log.Printf(suffix+msg, args...)
		case PANIC:
			log.Panicf(suffix+msg, args...)
		}
	}
}

func GetLoggerFromCtx(ctx context.Context) (*Logger, error) {
	logger, ok := ctx.Value(LOGGERCTX).(*Logger)
	if !ok {
		return nil, fmt.Errorf("logger not found in context")
	}
	return logger, nil
}

func CheckErrLog(err error, msg ...interface{}) {
	if err != nil {
		log.Println(fmt.Sprintf("%v", msg...) + err.Error())
	}
}

func CheckErrReturn(err error, msg ...interface{}) func() error {
	if err != nil {
		return func() error {
			return fmt.Errorf("%v: %w\n", fmt.Sprintf("%v", msg...), err)
		}
	}
	return func() error {
		return nil
	}
}

func CheckErrAll(err error, msg ...interface{}) func() error {
	CheckErrLog(err, msg)
	return CheckErrReturn(err, msg)
}

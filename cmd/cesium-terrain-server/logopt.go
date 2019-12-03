package main

import (
	"errors"
	"github.com/nmccready/cesium-terrain-server/log"
)

type LogOpt struct {
	Priority log.Priority
}

func NewLogOpt() *LogOpt {
	return &LogOpt{
		Priority: log.LOG_NOTICE,
	}
}

func (this *LogOpt) String() string {
	switch this.Priority {
	case log.LOG_CRIT:
		return "crit"
	case log.LOG_ERR:
		return "err"
	case log.LOG_NOTICE:
		return "notice"
	default:
		return "debug"
	}
}

func (this *LogOpt) Set(level string) error {
	switch level {
	case "crit":
		this.Priority = log.LOG_CRIT
	case "err":
		this.Priority = log.LOG_ERR
	case "notice":
		this.Priority = log.LOG_NOTICE
	case "debug":
		this.Priority = log.LOG_DEBUG
	default:
		return errors.New("choose one of crit, err, notice, debug")
	}
	return nil
}

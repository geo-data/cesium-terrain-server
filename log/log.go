package log

import (
	l "log"
	"os"
)

type Priority int

const (
	LOG_DEBUG Priority = iota
	LOG_NOTICE
	LOG_ERR
	LOG_CRIT
)

// Logger interface is satisfied by syslog.Writer
type Logger interface {
	Debug(m string) (err error)
	Notice(m string) (err error)
	Err(m string) (err error)
	Crit(m string) (err error)
}

var std Logger = New(l.New(os.Stderr, "", l.LstdFlags), LOG_NOTICE)

func Debug(m string) error {
	return std.Debug(m)
}

func Notice(m string) error {
	return std.Notice(m)
}

func Err(m string) error {
	return std.Err(m)
}

func Crit(m string) error {
	return std.Crit(m)
}

type logProxy struct {
	priority Priority
	log      *l.Logger
}

func New(logger *l.Logger, priority Priority) Logger {
	return &logProxy{
		log:      logger,
		priority: priority,
	}
}

func (this *logProxy) write(m string, p Priority) (err error) {
	if this.priority <= p {
		err = this.log.Output(2, m)
	}
	return
}

func (this *logProxy) Debug(m string) (err error) {
	return this.write("DEBUG: "+m, LOG_DEBUG)
}

func (this *logProxy) Notice(m string) (err error) {
	return this.write("NOTICE: "+m, LOG_NOTICE)
}

func (this *logProxy) Err(m string) (err error) {
	return this.write("ERROR: "+m, LOG_ERR)
}

func (this *logProxy) Crit(m string) (err error) {
	return this.write("CRITICAL: "+m, LOG_CRIT)
}

func SetLogger(logger Logger) {
	std = logger
}

func SetLog(log *l.Logger, priority Priority) {
	SetLogger(New(log, priority))
}

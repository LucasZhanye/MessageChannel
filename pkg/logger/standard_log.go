package logger

import (
	"fmt"
	stdlog "log"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

type defaultLog struct {
	log   *stdlog.Logger
	level Level
}

func InitDefaultLog(conf Config) (Log, error) {

	logLevel := map[string]Level{
		"debug": DEBUG,
		"info":  INFO,
		"warn":  WARN,
		"error": ERROR,
		"fatal": FATAL,
		"panic": PANIC,
	}

	level, ok := logLevel[conf.Level]
	if !ok {
		level = logLevel["info"]
	}

	log := stdlog.New(os.Stdout, "", stdlog.LstdFlags)
	log.SetFlags(stdlog.Ldate | stdlog.Lshortfile)

	return &defaultLog{
		log:   log,
		level: level,
	}, nil
}

func (l *defaultLog) printf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	l.log.Output(3, str)
}

func (l *defaultLog) println(args ...interface{}) {
	str := fmt.Sprintln(args...)
	l.log.Output(3, str)
}

func (l *defaultLog) Info(format string, args ...interface{}) {
	if INFO >= l.level {
		l.log.SetPrefix("[INFO] ")
		l.printf(format, args...)
	}
}

func (l *defaultLog) Debug(format string, args ...interface{}) {
	if DEBUG >= l.level {
		l.log.SetPrefix("[DEBUG] ")
		l.printf(format, args...)
	}
}

func (l *defaultLog) Error(format string, args ...interface{}) {
	if ERROR >= l.level {
		l.log.SetPrefix("[ERROR] ")
		l.printf(format, args...)
	}
}

func (l *defaultLog) Warn(format string, args ...interface{}) {
	if WARN >= l.level {
		l.log.SetPrefix("[WARN] ")
		l.printf(format, args...)
	}
}

func (l *defaultLog) Fatal(format string, args ...interface{}) {
	if FATAL >= l.level {
		l.log.SetPrefix("[FATAL] ")
		l.printf(format, args...)
	}
}

func (l *defaultLog) Panic(format string, args ...interface{}) {
	if PANIC >= l.level {
		l.log.SetPrefix("[PANIC] ")
		l.printf(format, args...)
	}
}

func (l defaultLog) Println(args ...interface{}) {
	if INFO >= l.level {
		l.log.SetPrefix("[INFO] ")
		l.println(args...)
	}
}

func (l defaultLog) Printf(format string, args ...interface{}) {
	if INFO >= l.level {
		l.log.SetPrefix("[INFO] ")
		l.printf(format, args...)
	}
}

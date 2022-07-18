package spoor

import (
	"fmt"
	"io"
	"time"
)

type Logger struct {
	callerSkip int
	level      Level
	Base
	files map[Level]*FileWriter
}

type Option func(l *Logger)

func CallerSkip(callerSkip int) Option {
	return func(l *Logger) {
		l.callerSkip = callerSkip
	}
}

func SetOutput(writer io.Writer) Option {
	return func(l *Logger) {
		l.Base.SetOutput(writer)
	}
}

func NewLogger(options ...Option) *Logger {
	l := &Logger{
		callerSkip: 0,
		Base:       Base{},
	}
	for _, option := range options {
		option(l)
	}
	return l
}

func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

func (l *Logger) createFiles(level Level) error {
	now := time.Now()

	for s := level; level >= DebugLog && l.files[level] == nil; s-- {
		fw := &FileWriter{
			level: s,
		}
		if err := fw.rotateFile(now); err != nil {
			return err
		}
		l.files[s] = fw
	}
	return nil
}

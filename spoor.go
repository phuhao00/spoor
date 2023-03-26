package spoor

import (
	"io"
	"log"
)

type Spoor struct {
	l        Logger
	cfgLevel Level
	prefix   string
	flag     int
}

type Option func(spoor *Spoor)

func WithFileWriter(writer *FileWriter) Option {
	return func(spoor *Spoor) {
		writer.level = spoor.cfgLevel
		spoor.l.SetOutput(writer)
	}
}

func WithConsoleWriter(writer io.Writer) Option {
	return func(spoor *Spoor) {
		spoor.l.SetOutput(writer)
	}
}

func NewSpoor(cfgLevel Level, prefix string, flag int, opts ...Option) *Spoor {
	logger := log.New(io.Discard, prefix, flag)
	s := &Spoor{
		l:        logger,
		cfgLevel: cfgLevel,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (l *Spoor) CheckLevel(level Level) bool {
	if level >= l.cfgLevel {
		return false
	}
	return true
}

package spoor

import (
	"io"
	"log"
)

type Spoor struct {
	Logger
	cfgLevel Level
	prefix   string
	flag     int
}

type Option func(spoor *Spoor)

func WithFileWriter(writer *FileWriter) Option {
	return func(spoor *Spoor) {
		writer.level = spoor.cfgLevel
		spoor.SetOutput(writer)
	}
}

func WithConsoleWriter(writer io.Writer) Option {
	return func(spoor *Spoor) {
		spoor.SetOutput(writer)
	}
}

func NewSpoor(cfgLevel Level, prefix string, flag int, opts ...Option) *Spoor {
	logger := log.New(io.Discard, prefix, flag)
	s := &Spoor{
		Logger:   logger,
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

type LoggingSetting struct {
	Dir          string
	Level        int
	Prefix       string
	WriterOption Option
}

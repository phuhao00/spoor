package logger

import (
	"fmt"
	"github.com/phuhao00/spoor"
	"log"
	"sync"
)

var (
	sp             *spoor.Spoor
	onceInitLogger sync.Once
)

func GetLogger() *spoor.Spoor {
	return sp
}

type LoggingSetting struct {
	Dir          string
	Level        int
	Prefix       string
	WriterOption spoor.Option
}

func SetLogging(setting *LoggingSetting) {
	onceInitLogger.Do(func() {
		var opt spoor.Option
		if setting.WriterOption == nil {
			fileWriter := spoor.NewFileWriter(setting.Dir, 0, 0, 0)
			opt = spoor.WithFileWriter(fileWriter)
		} else {
			opt = setting.WriterOption
		}
		l := spoor.NewSpoor(spoor.Level(setting.Level), setting.Prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, opt)
		sp = l
	})
}

// Debug Log line format: [IWEF]mmdd hh:mm:sLogger.uuuuuu threadid file:line] msg
func Debug(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.DEBUG) {
		return
	}
	sp.Output(2, fmt.Sprintf(spoor.DEBUG.String()+" "+f, args...))
}

func Error(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.ERROR) {
		return
	}
	sp.Output(2, fmt.Sprintf(spoor.ERROR.String()+" "+f, args...))
}

func Info(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.INFO) {
		return
	}
	sp.Output(2, fmt.Sprintf(spoor.INFO.String()+" "+f, args...))
}

func Warn(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.WARN) {
		return
	}
	sp.Output(2, fmt.Sprintf(spoor.WARN.String()+" "+f, args...))
}

func Fatal(f string, args ...interface{}) {
	if sp.CheckLevel(spoor.FATAL) {
		return
	}
	sp.Output(2, fmt.Sprintf(spoor.FATAL.String()+" "+f, args...))
}

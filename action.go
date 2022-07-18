package spoor

func (l *Logger) DebugF(format string, args ...interface{}) {
	err := l.Output(l.callerSkip, getMessage(format, args))
	if err != nil {

	}
}

func (l *Logger) InfoF(format string, args ...interface{}) {

}

func (l *Logger) ErrorF(format string, args ...interface{}) {

}

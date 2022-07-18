package spoor

type Level int

const (
	DebugLog Level = iota
	InfoLog
	WarningLog
	ErrorLog
	FatalLog
)

func (l Level) ToString() string {
	switch l {
	case DebugLog:
		return "DebugLog"
	case InfoLog:
		return "InfoLog"
	case WarningLog:
		return "WarningLog"
	case ErrorLog:
		return "ErrorLog"
	case FatalLog:
		return "FatalLog"
	}
	return ""
}

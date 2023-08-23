package helpers

import "log"

type LogLevel string

const (
	Info  = "INFO"
	Warn  = "WARN"
	Error = "ERROR"
)

func logLevelString(l LogLevel) string {
	switch l {
	case Info:
		return "[INFO]"
	case Warn:
		return "[WARN]"
	case Error:
		return "[ERROR]"
	default:
		return "[INFO]"
	}
}

func Log(l LogLevel, e error) {
	level := logLevelString(l)
	log.Printf("%s: %s", level, e)
}

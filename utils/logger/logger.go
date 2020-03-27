package logger

import "time"

const (
	INFO = iota
	DBG
	WARN
	ERR
)

func Log(Level int, msg string) {
	switch Level {
	case INFO:
		logInfo(msg)
	case DBG:
		logDebug(msg)
	case WARN:
		logWarn(msg)
	case ERR:
		logErr(msg)
	default:
		logDebug(msg)
	}
}

func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func logInfo(msg string) {
	println(getTime() + " [INFO]: " + msg)
}

func logErr(msg string) {
	println(getTime() + " [ERROR]: " + msg)
}

func logWarn(msg string) {
	println(getTime() + " [WARNING]: " + msg)
}

func logDebug(msg string) {
	println(getTime() + " [DEBUG]: " + msg)
}

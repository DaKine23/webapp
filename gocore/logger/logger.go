package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

// 		Logging according to RFC 5424

// lvl     meaning
//______________________________________________
// 0       Emergency: system is unusable
// 1       Alert: action must be taken immediately
// 2       Critical: critical conditions
// 3       Error: error conditions
// 4       Warning: warning conditions
// 5       Notice: normal but significant condition
// 6       Informational: informational messages
// 7       Debug: debug-level messages

type LogLevel int

var ll LogLevel = Informational

//Init sets global LogLevel to a new value any Loglevel that is equal or lower will be logged
// lvl     meaning
//______________________________________________
// 0       Emergency: system is unusable
// 1       Alert: action must be taken immediately
// 2       Critical: critical conditions
// 3       Error: error conditions
// 4       Warning: warning conditions
// 5       Notice: normal but significant condition
// 6       Informational: informational messages
// 7       Debug: debug-level messages
func Init(logLevel LogLevel) {
	ll = logLevel
}

const (
	NoLog         LogLevel = -1   //turn off logs if used for Init and NoOp if used for Log
	Emergency     LogLevel = iota //Emergency: system is unusable
	Alert         LogLevel = iota //Alert: action must be taken immediately
	Critical      LogLevel = iota //Critical: critical conditions
	Error         LogLevel = iota //Error: error conditions
	Warning       LogLevel = iota //Warning: warning conditions
	Notice        LogLevel = iota //Notice: normal but significant condition
	Informational LogLevel = iota //Informational: informational messages
	Debug         LogLevel = iota //Debug: debug-level messages

	defaultSkip = 3 // skipping caller
)

// LogEmergency logs a message with EMERGENCY log level
func LogEmergency(flowID string, format string, values ...interface{}) {
	logEmergency(defaultSkip, flowID, format, values...)
}
func logEmergency(skip int, flowID string, format string, values ...interface{}) {
	if !(ll < Emergency) {
		logAbstract(skip, format, prepareValues(flowID, "EMERGENCY", values)...)
	}
}

// LogAlert logs a message with CRITICAL log level
func LogAlert(flowID string, format string, values ...interface{}) {
	logAlert(defaultSkip, flowID, format, values...)
}
func logAlert(skip int, flowID string, format string, values ...interface{}) {
	if !(ll < Alert) {
		logAbstract(skip, format, prepareValues(flowID, "ALERT", values)...)
	}
}

// LogCritical logs a message with CRITICAL log level
func LogCritical(flowID string, format string, values ...interface{}) {
	logCritical(defaultSkip, flowID, format, values...)
}
func logCritical(skip int, flowID string, format string, values ...interface{}) {
	if !(ll < Critical) {
		logAbstract(skip, format, prepareValues(flowID, "CRITICAL", values)...)
	}
}

// LogError logs a message with ERROR log level
func LogError(flowID string, format string, values ...interface{}) {
	logError(defaultSkip, flowID, format, values...)
}
func logError(skip int, flowID string, format string, values ...interface{}) {
	if !(ll < Error) {
		logAbstract(skip, format, prepareValues(flowID, "ERROR", values)...)
	}
}

// LogWarning logs a message with WARNING log level
func LogWarning(flowID string, format string, values ...interface{}) {
	logWarning(defaultSkip, flowID, format, values...)
}
func logWarning(skip int, flowID string, format string, values ...interface{}) {
	if !(ll < Warning) {
		logAbstract(skip, format, prepareValues(flowID, "WARNING", values)...)
	}
}

// LogNotice logs a message with ERROR log level
func LogNotice(flowID string, format string, values ...interface{}) {
	logNotice(defaultSkip, flowID, format, values...)
}
func logNotice(skip int, flowID string, format string, values ...interface{}) {
	if !(ll < Notice) {
		logAbstract(skip, format, prepareValues(flowID, "NOTICE", values)...)
	}
}

// LogInfo logs a message with INFO log level
func LogInfo(flowID string, format string, values ...interface{}) {
	logInfo(defaultSkip, flowID, format, values...)
}
func logInfo(skip int, flowID string, format string, values ...interface{}) {
	if !(ll < Informational) {
		logAbstract(skip, format, prepareValues(flowID, "INFO", values)...)
	}
}

// LogDebug logs a message with DEBUG log level
func LogDebug(flowID string, format string, values ...interface{}) {
	logDebug(defaultSkip, flowID, format, values...)
}
func logDebug(skip int, flowID string, format string, values ...interface{}) {
	if !(ll < Debug) {
		logAbstract(skip, format, prepareValues(flowID, "DEBUG", values)...)
	}
}

// Combine flowID and logLevel with other values
func prepareValues(flowID string, logLevel string, values []interface{}) []interface{} {
	return append([]interface{}{flowID, logLevel}, values...)
}

// LogS logs a message with given LogLevel and additional skip (relative to the default behavior
func LogS(lvl LogLevel, skip int, flowID string, format string, values ...interface{}) {

	switch lvl {
	case Emergency:
		logEmergency(defaultSkip+skip, flowID, format, values...)
	case Alert:
		logAlert(defaultSkip+skip, flowID, format, values...)
	case Critical:
		logCritical(defaultSkip+skip, flowID, format, values...)
	case Error:
		logError(defaultSkip+skip, flowID, format, values...)
	case Warning:
		logWarning(defaultSkip+skip, flowID, format, values...)
	case Notice:
		logNotice(defaultSkip+skip, flowID, format, values...)
	case Informational:
		logInfo(defaultSkip+skip, flowID, format, values...)
	case Debug:
		logDebug(defaultSkip+skip, flowID, format, values...)
	}

}

// Log logs a message with given LogLevel
func Log(lvl LogLevel, flowID string, format string, values ...interface{}) {

	LogS(lvl, 1, flowID, format, values...)

}

// Logging according to RFC 5424
func logAbstract(skip int, format string, values ...interface{}) {

	_, fn, line, _ := runtime.Caller(skip)
	log.Printf("[%s] [%s] "+fmt.Sprintf("[%s:%d] ", removeFullPath(fn), line)+format, values...)

}

func removeFullPath(fn string) string {

	//REPO_PROVIDER allows to use animal in all kinds of repos e.g. github.com or bitbucket
	repoProvider := os.Getenv("REPO_PROVIDER")

	if len(repoProvider) == 0 {
		repoProvider = "github.bus.zalan.do"
	}

	return removeFromStartToInputIfFound(fn, repoProvider)
}

func removeFromStartToInputIfFound(input, remove string) string {

	index := strings.Index(input, remove)
	if index == -1 {
		index = 0
	} else {
		index += len(remove) + 1
	}
	return input[index:]

}

package sdk

import (
	"fmt"
	"log"
	"os"
)

var (
	gLogger   ILogger  = log.New(os.Stdout, "[eventbus-sdk] ", log.Ldate|log.Ltime)
	gLogLevel LogLevel = LogLevelInfo
)

type ILogger interface {
	Output(maxdepth int, s string) error
}

//LogLevel specifies the severity of a given log message
type LogLevel int

//Log levels
const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// String returns the string form for a given LogLevel
func (lvl LogLevel) String() string {
	switch lvl {
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	}
	return "DEBUG"
}

// SetLogger sets global scope log options
// These global log options can be overridden by the Producer.SetLogger or Consumer.SetLogger.
func SetLogger(l ILogger, lvl LogLevel) {
	gLogger = l
	gLogLevel = lvl
}

// Output debug level log
func Debug(format string, v ...interface{}) {
	if gLogger != nil {
		if gLogLevel <= LogLevelDebug {
			s := "[DEBUG] " + fmt.Sprintf(format, v...)
			gLogger.Output(2, s)
		}
	}
}

// Output info level log
func Info(format string, v ...interface{}) {
	if gLogger != nil {
		if gLogLevel <= LogLevelInfo {
			s := "[INFO] " + fmt.Sprintf(format, v...)
			gLogger.Output(2, s)
		}
	}
}

// Output warn level log
func Warn(format string, v ...interface{}) {
	if gLogger != nil {
		if gLogLevel <= LogLevelWarn {
			s := "[WARN] " + fmt.Sprintf(format, v...)
			gLogger.Output(2, s)
		}
	}
}

// Output error level log
func Error(format string, v ...interface{}) {
	if gLogger != nil {
		if gLogLevel <= LogLevelError {
			s := "[ERROR] " + fmt.Sprintf(format, v...)
			gLogger.Output(2, s)
		}
	}
}

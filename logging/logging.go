package logging

import (
	"github.com/boreq/blogs/config"
	"log"
	"os"
)

// Logger defines methods used for logging in a normal mode and a debug mode.
// Debug mode log messages are displayed only if a proper environment variable
// with the name stored in DebugEnvVar is set.
type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
}

var debug *bool
var loggers map[string]Logger

func init() {
	loggers = make(map[string]Logger)
	debug = &config.Config.Debug
}

// GetLogger creates a new logger or returns an already existing logger created
// with the given name using this method.
func GetLogger(name string) Logger {
	if _, ok := loggers[name]; !ok {
		loggers[name] = &logger{log.New(os.Stdout, name+": ", 0)}
	}
	return loggers[name]
}

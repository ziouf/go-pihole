package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

// Errors
var (
	ErrLogLevelNotFound = errors.New(`Log level not found`)
)

// Log levels
var (
	VERBOSE = lvl{id: 0, label: "VERBOSE"}
	DEBUG   = lvl{id: 1, label: "DEBUG"}
	INFO    = lvl{id: 2, label: "INFO"}
	ERROR   = lvl{id: 3, label: "ERROR"}
)
var lvls = []lvl{VERBOSE, DEBUG, INFO, ERROR}

func getLevel(l string) lvl {
	for _, lvl := range lvls {
		if lvl.label == l {
			return lvl
		}
	}
	return INFO
}

var logger logging
var loggerMap = make(map[int]*log.Logger, 0)

// Init logger configuration
func Init(file string, lvl string) {
	logger = logging{path: file, level: getLevel(lvl)}
}

type logging struct {
	path  string
	level lvl
}

type lvl struct {
	id    int
	label string
}

func (lvl lvl) String() string {
	return lvl.label
}

func (lvl lvl) Equals(s string) bool {
	return strings.ToUpper(lvl.label) == strings.ToUpper(s)
}

func get(lvl lvl) *log.Logger {
	l, ok := loggerMap[lvl.id]
	if !ok {
		l = log.New(os.Stderr, fmt.Sprintf("[%s]", strings.ToUpper(lvl.label)), log.LstdFlags)
		if len(logger.path) == 0 {
			l.SetFlags(0)
			l.SetOutput(os.Stderr)
		} else {
			l.SetFlags(log.LstdFlags)
			l.SetOutput(newRotateWriter(path.Join(logger.path, lvl.label+".log")))
		}
		loggerMap[lvl.id] = l
	}
	return l
}

// Verbose return verbose logger
func Verbose() *log.Logger {
	return get(VERBOSE)
}

// Debug return debug logger
func Debug() *log.Logger {
	return get(DEBUG)
}

// Info return info logger
func Info() *log.Logger {
	return get(INFO)
}

// Error return error logger
func Error() *log.Logger {
	return get(ERROR)
}

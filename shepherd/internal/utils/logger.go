package utils

import (
	"fmt"
	"io"
	"log"
)

// Logger - wrapper for log.Logger instances, providing both
// info and error loggers.
type Logger struct {
	Info *log.Logger
	Err  *log.Logger
}

// NewLogger - returns new Logger instance.
func NewLogger(mainTag string, infoOut, errOut io.Writer, extraTags ...string) *Logger {
	extraPrefix := ""

	for _, tag := range extraTags {
		extraPrefix += fmt.Sprintf("[%s]", tag)
	}

	extraPrefix += ": "

	return &Logger{
		Info: log.New(infoOut, fmt.Sprintf("[%s][Info]%s", mainTag, extraPrefix), log.Ldate|log.Ltime|log.Lmsgprefix),
		Err:  log.New(errOut, fmt.Sprintf("[%s][Error]%s", mainTag, extraPrefix), log.Ldate|log.Ltime|log.Lmsgprefix),
	}
}

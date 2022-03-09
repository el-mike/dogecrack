package common

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

	infoPrefix string
	errPrefix  string
}

// NewLogger - returns new Logger instance.
func NewLogger(mainTag string, infoOut, errOut io.Writer, extraTags ...string) *Logger {
	extraPrefix := ""

	for _, tag := range extraTags {
		extraPrefix += fmt.Sprintf("[%s]", tag)
	}

	extraPrefix += ": "

	infoPrefix := fmt.Sprintf("[%s][Info]%s", mainTag, extraPrefix)
	errPrefix := fmt.Sprintf("[%s][Error]%s", mainTag, extraPrefix)

	return &Logger{
		Info:       log.New(infoOut, infoPrefix, log.Ldate|log.Ltime|log.Lmsgprefix),
		Err:        log.New(errOut, errPrefix, log.Ldate|log.Ltime|log.Lmsgprefix),
		infoPrefix: infoPrefix,
		errPrefix:  errPrefix,
	}
}

// DecorateErr - returns given error message with it's prefix.
func (l *Logger) DecorateErr(err error) error {
	return fmt.Errorf("%s%s", l.errPrefix, err)
}

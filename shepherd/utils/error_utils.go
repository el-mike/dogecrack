package utils

import (
	"fmt"
	"log"
	"runtime"
)

func getRuntimeInfo(err error) string {
	callStack := ""
	// We need to skip couple of Callers to get the actual info.
	// 0 - LogRuntimeInfo
	// 1 - caller of LogRuntimeInfo
	// 2 - the actual error location
	stackPos := 2

	_, fn, line, ok := runtime.Caller(stackPos)
	for ok {
		callStack += fmt.Sprintf("%s:%d\n", fn, line)

		stackPos += 1
		_, fn, line, ok = runtime.Caller(stackPos)
	}

	return fmt.Sprintf("%v\ncall stack:\n%v", err.Error(), callStack)
}

// WithStackTrace - logs the error and its stack trace with given logger,
// and returns the error back.
func WithStackTrace(logger *log.Logger, err error) error {
	logger.Print(getRuntimeInfo(err))

	return err
}

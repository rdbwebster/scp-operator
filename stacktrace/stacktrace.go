package stacktrace

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type StackTrace struct {
	msg  string
	path string
}

// New function constructs a new `StackTrace` struct by using given panic
// message, absolute path of the caller file and the line number.
func New(msg string) *StackTrace {
	_, file, line, _ := runtime.Caller(1)
	p, _ := os.Getwd()

	return &StackTrace{
		msg:  msg,
		path: fmt.Sprintf("%s:%d", strings.TrimPrefix(file, p), line),
	}
}

func (s *StackTrace) Message() string {
	return s.msg
}

func (s *StackTrace) Path() string {
	return s.path
}

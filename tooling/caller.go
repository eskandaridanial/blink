package tooling

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// func 'Caller' returns a short 'file.go:line' string for the caller at the given skip depth.
// skip=0 means 'Caller' itself, skip=1 is the immediate caller, and so on.
func Caller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

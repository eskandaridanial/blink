package tooling

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func Caller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

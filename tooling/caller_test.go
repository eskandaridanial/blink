package tooling

import (
	"strings"
	"testing"
)

func TestCaller(t *testing.T) {
	call := Caller(1)
	if !strings.HasPrefix(call, "caller_test.go:") {
		t.Fatalf("unexpected caller format: got %s", call)
	}
}

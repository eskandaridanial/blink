package formatters

import (
	"fmt"
	"time"

	"github.com/eskandaridanial/blink/models"
)

type TextFormatter struct{}

func (f *TextFormatter) Format(r models.Record) []byte {
	s := fmt.Sprintf(
		"[%s] %s %s %s - %s",
		r.Level.String(),
		r.Timestamp.Format(time.RFC3339),
		r.Caller,
		r.ReferenceId,
		r.Message,
	)

	for _, field := range r.Fields {
		s += fmt.Sprintf(" %s=%v", field.Key, field.Value)
	}

	return []byte(s)
}

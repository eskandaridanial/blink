package formatters

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/eskandaridanial/blink/models"
)

var textPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

type TextFormatter struct{}

func (f *TextFormatter) Format(r models.Record) []byte {
	b := textPool.Get().(*bytes.Buffer)
	b.Reset()
	defer textPool.Put(b)

	b.WriteString(r.Timestamp.Format(time.RFC3339))
	b.WriteByte(' ')
	b.WriteString(r.Level.String())
	b.WriteByte(' ')
	b.WriteString(r.ReferenceId)
	b.WriteByte(' ')
	b.WriteString(r.Caller)
	b.WriteString(" [")

	for _, field := range r.Fields {
		b.WriteString(" ")
		b.WriteString(field.Key)
		b.WriteString("=")
		switch val := field.Value.(type) {
		case string:
			b.WriteString(val)
		case fmt.Stringer:
			b.WriteString(val.String())
		default:
			b.WriteString(fmt.Sprintf("%v", val))
		}
	}
	b.WriteString(" ]")
	b.WriteString(" - ")
	b.WriteString(r.Message)
	b.WriteByte('\n')
	return append([]byte(nil), b.String()...)
}

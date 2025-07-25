package formatters

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"

	"github.com/eskandaridanial/blink/models"
)

var jsonPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

type JsonFormatter struct{}

func (f *JsonFormatter) Format(r models.Record) []byte {
	buf := jsonPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer jsonPool.Put(buf)

	payload := make(map[string]any, 6)
	payload["timestamp"] = r.Timestamp.Format(time.RFC3339)
	payload["level"] = r.Level.String()
	payload["referenceId"] = r.ReferenceId
	payload["caller"] = r.Caller
	payload["message"] = r.Message

	if len(r.Fields) > 0 {
		fields := make(map[string]any, len(r.Fields))
		for _, f := range r.Fields {
			fields[f.Key] = f.Value
		}
		payload["fields"] = fields
	}

	_ = json.NewEncoder(buf).Encode(payload)
	return append([]byte(nil), buf.Bytes()...)
}

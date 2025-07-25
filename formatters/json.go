package formatters

import (
	"encoding/json"
	"time"

	"github.com/eskandaridanial/blink/models"
)

type JsonFormatter struct{}

func (f *JsonFormatter) Format(record models.Record) []byte {
	out := make(map[string]any, 6)
	out["level"] = record.Level.String()
	out["timestamp"] = record.Timestamp.Format(time.RFC3339)
	out["caller"] = record.Caller
	out["message"] = record.Message
	out["referenceId"] = record.ReferenceId

	if len(record.Fields) > 0 {
		fields := make(map[string]any, len(record.Fields))
		for _, field := range record.Fields {
			fields[field.Key] = field.Value
		}
		out["fields"] = fields
	}

	result, _ := json.Marshal(out)
	return append(result, '\n')
}

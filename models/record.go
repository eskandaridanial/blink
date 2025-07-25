package models

import (
	"time"
)

type Record struct {
	Level       Level
	Message     string
	Caller      string
	ReferenceId string
	Fields      []Field
	Timestamp   time.Time
}

type RecordBuilder struct {
	record Record
}

func NewRecordBuilder() RecordBuilder {
	return RecordBuilder{
		record: Record{
			Timestamp: time.Now(),
		},
	}
}

func (b RecordBuilder) Level(level Level) RecordBuilder {
	b.record.Level = level
	return b
}

func (b RecordBuilder) Message(msg string) RecordBuilder {
	b.record.Message = msg
	return b
}

func (b RecordBuilder) ReferenceId(id string) RecordBuilder {
	b.record.ReferenceId = id
	return b
}

func (b RecordBuilder) Caller(caller string) RecordBuilder {
	b.record.Caller = caller
	return b
}

func (b RecordBuilder) Fields(fields []Field) RecordBuilder {
	copied := append([]Field(nil), fields...)
	b.record.Fields = copied
	return b
}

func (b RecordBuilder) Timestamp(t time.Time) RecordBuilder {
	b.record.Timestamp = t
	return b
}

func (b RecordBuilder) Build() Record {
	return b.record
}

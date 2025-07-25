package logging

import (
	"context"

	"github.com/eskandaridanial/blink/handlers"
	"github.com/eskandaridanial/blink/hooks"
	"github.com/eskandaridanial/blink/models"
	"github.com/eskandaridanial/blink/tooling"
)

type Logger struct {
	Level  models.Level
	Fields []models.Field

	HookMux    hooks.HookMux
	HandlerMux handlers.HandlerMux
}

func NewLogger(opts ...Option) *Logger {
	l := &Logger{}
	for _, o := range opts {
		o(l)
	}
	return l
}

func (l *Logger) WithField(field models.Field) *Logger {
	newFields := append([]models.Field(nil), l.Fields...)
	newFields = append(newFields, field)

	return &Logger{
		Level:      l.Level,
		Fields:     newFields,
		HookMux:    l.HookMux,
		HandlerMux: l.HandlerMux,
	}
}

func (l *Logger) WithFields(fields []models.Field) *Logger {
	newFields := append([]models.Field(nil), l.Fields...)
	newFields = append(newFields, fields...)

	return &Logger{
		Level:      l.Level,
		Fields:     newFields,
		HookMux:    l.HookMux,
		HandlerMux: l.HandlerMux,
	}
}

func (l *Logger) Log(ctx context.Context, level models.Level, message string, fields ...models.Field) {
	if level < l.Level {
		return
	}

	combined := append([]models.Field(nil), l.Fields...)
	combined = append(combined, fields...)

	record := models.NewRecordBuilder().
		Level(level).
		Message(message).
		Fields(combined).
		Caller(tooling.Caller(3)).
		Build()

	l.HookMux.Apply(ctx, record)
	l.HandlerMux.Apply(ctx, record)
}

func (l *Logger) Debug(message string, fields ...models.Field) {
	l.Log(context.Background(), models.Debug, message, fields...)
}

func (l *Logger) Info(message string, fields ...models.Field) {
	l.Log(context.Background(), models.Info, message, fields...)
}

func (l *Logger) Warn(message string, fields ...models.Field) {
	l.Log(context.Background(), models.Warn, message, fields...)
}

func (l *Logger) Error(message string, fields ...models.Field) {
	l.Log(context.Background(), models.Error, message, fields...)
}

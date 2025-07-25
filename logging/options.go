package logging

import (
	"github.com/eskandaridanial/blink/handlers"
	"github.com/eskandaridanial/blink/hooks"
	"github.com/eskandaridanial/blink/models"
)

type Option func(*Logger)

func WithLevel(level models.Level) Option {
	return func(l *Logger) {
		l.Level = level
	}
}

func WithField(field models.Field) Option {
	return func(l *Logger) {
		l.Fields = append([]models.Field(nil), field)
	}
}

func WithFields(fields ...models.Field) Option {
	return func(l *Logger) {
		l.Fields = append([]models.Field(nil), fields...)
	}
}

func WithHandler(handler any) Option {
	return func(l *Logger) {
		l.HandlerMux = *handlers.NewHandlerMuxBuilder().WithHandler(handler).Build()
	}
}

func WithHook(hook any) Option {
	return func(l *Logger) {
		l.HookMux = *hooks.NewHookMuxBuilder().WithHook(hook).Build()
	}
}

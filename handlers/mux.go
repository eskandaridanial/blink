package handlers

import (
	"context"

	"github.com/eskandaridanial/blink/models"
)

type HandlerMux struct {
	handlers []Handler
}

type HandlerMuxBuilder struct {
	handlers []Handler
}

func NewHandlerMuxBuilder() *HandlerMuxBuilder {
	return &HandlerMuxBuilder{}
}

func (b *HandlerMuxBuilder) WithHandler(h any) *HandlerMuxBuilder {
	if handler, ok := h.(Handler); ok {
		b.handlers = append(b.handlers, handler)
	}
	return b
}

func (b *HandlerMuxBuilder) Build() *HandlerMux {
	return &HandlerMux{handlers: b.handlers}
}

func (m *HandlerMux) Apply(ctx context.Context, r models.Record) {
	for _, h := range m.handlers {
		h.Handle(ctx, r)
	}
}

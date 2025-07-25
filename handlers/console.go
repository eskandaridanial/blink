package handlers

import (
	"context"
	"os"

	"github.com/eskandaridanial/blink/formatters"
	"github.com/eskandaridanial/blink/models"
)

type ConsoleHandler struct {
	Formatter formatters.Formatter
}

func (h *ConsoleHandler) Handle(ctx context.Context, r models.Record) {
	formatted := h.Formatter.Format(r)
	_, _ = os.Stdout.Write(formatted)
}

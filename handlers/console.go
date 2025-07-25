package handlers

import (
	"bufio"
	"context"
	"os"
	"sync"

	"github.com/eskandaridanial/blink/formatters"
	"github.com/eskandaridanial/blink/models"
)

type ConsoleHandler struct {
	formatter formatters.Formatter
	writer    *bufio.Writer
	mu        sync.Mutex
}

func NewConsoleHandler(formatter formatters.Formatter) *ConsoleHandler {
	return &ConsoleHandler{
		formatter: formatter,
		writer:    bufio.NewWriterSize(os.Stdout, 4096),
	}
}

func (h *ConsoleHandler) Handle(ctx context.Context, r models.Record) {
	formatted := h.formatter.Format(r)
	h.mu.Lock()
	_, _ = h.writer.Write(formatted)
	_ = h.writer.Flush()
	h.mu.Unlock()
}

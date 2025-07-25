package handlers

import (
	"context"

	"github.com/eskandaridanial/blink/models"
)

type Handler interface {
	Handle(ctx context.Context, r models.Record)
}

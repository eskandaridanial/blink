package hooks

import (
	"context"
	"fmt"

	"github.com/eskandaridanial/blink/models"
)

type DefaultHook struct{}

func (h *DefaultHook) OnAll(ctx context.Context, r models.Record) {
	fmt.Printf(
		"hook triggered: OnAll | level=%s | caller=%s | ref=%s | msg=%q\n",
		r.Level.String(),
		r.Caller,
		r.ReferenceId,
		r.Message,
	)
}

func (h *DefaultHook) OnDebug(ctx context.Context, r models.Record) {
	fmt.Printf(
		"hook triggered: OnDebug | level=%s | caller=%s | ref=%s | msg=%q\n",
		r.Level.String(),
		r.Caller,
		r.ReferenceId,
		r.Message,
	)
}

func (h *DefaultHook) OnInfo(ctx context.Context, r models.Record) {
	fmt.Printf(
		"hook triggered: OnInfo | level=%s | caller=%s | ref=%s | msg=%q\n",
		r.Level.String(),
		r.Caller,
		r.ReferenceId,
		r.Message,
	)
}

func (h *DefaultHook) OnWarn(ctx context.Context, r models.Record) {
	fmt.Printf(
		"hook triggered: OnWarn | level=%s | caller=%s | ref=%s | msg=%q\n",
		r.Level.String(),
		r.Caller,
		r.ReferenceId,
		r.Message,
	)
}

func (h *DefaultHook) OnError(ctx context.Context, r models.Record) {
	fmt.Printf(
		"hook triggered: OnError | level=%s | caller=%s | ref=%s | msg=%q\n",
		r.Level.String(),
		r.Caller,
		r.ReferenceId,
		r.Message,
	)
}

package hooks

import (
	"context"

	"github.com/eskandaridanial/blink/models"
)

type OnAllHook interface {
	OnAll(ctx context.Context, r models.Record)
}

type OnDebugHook interface {
	OnDebug(ctx context.Context, r models.Record)
}

type OnInfoHook interface {
	OnInfo(ctx context.Context, r models.Record)
}

type OnWarnHook interface {
	OnWarn(ctx context.Context, r models.Record)
}

type OnErrorHook interface {
	OnError(ctx context.Context, r models.Record)
}

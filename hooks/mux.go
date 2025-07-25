package hooks

import (
	"context"

	"github.com/eskandaridanial/blink/models"
)

type HookMux struct {
	allHooks   []OnAllHook
	debugHooks []OnDebugHook
	infoHooks  []OnInfoHook
	warnHooks  []OnWarnHook
	errorHooks []OnErrorHook
}

type HookMuxBuilder struct {
	allHooks   []OnAllHook
	debugHooks []OnDebugHook
	infoHooks  []OnInfoHook
	warnHooks  []OnWarnHook
	errorHooks []OnErrorHook
}

func NewHookMuxBuilder() *HookMuxBuilder {
	return &HookMuxBuilder{}
}

func (b *HookMuxBuilder) WithHook(h any) *HookMuxBuilder {
	if hook, ok := h.(OnAllHook); ok {
		b.allHooks = append(b.allHooks, hook)
	}
	if hook, ok := h.(OnDebugHook); ok {
		b.debugHooks = append(b.debugHooks, hook)
	}
	if hook, ok := h.(OnInfoHook); ok {
		b.infoHooks = append(b.infoHooks, hook)
	}
	if hook, ok := h.(OnWarnHook); ok {
		b.warnHooks = append(b.warnHooks, hook)
	}
	if hook, ok := h.(OnErrorHook); ok {
		b.errorHooks = append(b.errorHooks, hook)
	}
	return b
}

func (b *HookMuxBuilder) Build() *HookMux {
	return &HookMux{
		allHooks:   b.allHooks,
		debugHooks: b.debugHooks,
		infoHooks:  b.infoHooks,
		warnHooks:  b.warnHooks,
		errorHooks: b.errorHooks,
	}
}

func (m *HookMux) Apply(ctx context.Context, r models.Record) {
	for _, h := range m.allHooks {
		h.OnAll(ctx, r)
	}

	switch r.Level {
	case models.Debug:
		for _, h := range m.debugHooks {
			h.OnDebug(ctx, r)
		}
	case models.Info:
		for _, h := range m.infoHooks {
			h.OnInfo(ctx, r)
		}
	case models.Warn:
		for _, h := range m.warnHooks {
			h.OnWarn(ctx, r)
		}
	case models.Error:
		for _, h := range m.errorHooks {
			h.OnError(ctx, r)
		}
	}
}

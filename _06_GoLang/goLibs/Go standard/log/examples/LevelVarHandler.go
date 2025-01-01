package examples

import (
	"context"
	. "log/slog"
)

type LevelVarHandler struct {
	level   *LevelVar
	handler Handler
}

func NewLevelVarHandler(h Handler, level *LevelVar) *LevelVarHandler {
	if lh, ok := h.(*LevelVarHandler); ok {
		h = lh.Handler()
	}
	return &LevelVarHandler{level, h}
}

func (h *LevelVarHandler) Enabled(_ context.Context, level Level) bool {
	return level >= h.level.Level()
}

func (h *LevelVarHandler) Handle(ctx context.Context, r Record) error {
	return h.handler.Handle(ctx, r)
}

func (h *LevelVarHandler) WithAttrs(attrs []Attr) Handler {
	return NewLevelVarHandler(h.handler.WithAttrs(attrs), h.level)
}

func (h *LevelVarHandler) WithGroup(name string) Handler {
	return NewLevelVarHandler(h.handler.WithGroup(name), h.level)
}

func (h *LevelVarHandler) Handler() Handler {
	return h.handler
}

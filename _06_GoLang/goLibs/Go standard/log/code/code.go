package gostd

import (
	"context"
	"fmt"
	"log/slog"
)

func _logCase(_case string) {
	_logfln("case : %s", _case)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("LOG ERROR: \n%s", err)
}

func _log(s any) {
	fmt.Println(s)
}
func _logfln(format string, args ...any) {
	fmt.Printf(format+"\n", args...)

}

// wrap slog default handler before call SetDefault
type wrappingHandler struct {
	h slog.Handler
	l slog.Level
}

func (h *wrappingHandler) Set(level slog.Level) {
	h.l = level
}
func (h *wrappingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.l
}
func (h *wrappingHandler) WithGroup(name string) slog.Handler              { return h.h.WithGroup(name) }
func (h *wrappingHandler) WithAttrs(as []slog.Attr) slog.Handler           { return h.h.WithAttrs(as) }
func (h *wrappingHandler) Handle(ctx context.Context, r slog.Record) error { return h.h.Handle(ctx, r) }

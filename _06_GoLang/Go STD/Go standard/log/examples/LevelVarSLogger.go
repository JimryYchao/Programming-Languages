package examples

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

var ignorePC = false

// 不受 SetLogLoggerLevel  的影响
type LevelVarSLogger struct {
	level   *slog.LevelVar
	handler slog.Handler
}

func NewLevelVarSLogger(h slog.Handler, level slog.Level) *LevelVarSLogger {
	leVar := &slog.LevelVar{}
	leVar.Set(level)
	return &LevelVarSLogger{leVar, h}
}

func (l *LevelVarSLogger) SetLevel(level slog.Level) {
	l.level.Set(level)
}
func (l *LevelVarSLogger) SetHandler(h slog.Handler) {
	l.handler = h
}
func (l *LevelVarSLogger) Enabled(ctx context.Context, level slog.Level) bool {
	if level < l.level.Level() {
		return false
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return l.Handler().Enabled(ctx, level)
}

func (l *LevelVarSLogger) Debug(msg string, args ...any) {
	l.log(context.Background(), slog.LevelDebug, msg, args...)
}
func (l *LevelVarSLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelDebug, msg, args...)
}

func (l *LevelVarSLogger) Error(msg string, args ...any) {
	l.log(context.Background(), slog.LevelError, msg, args...)
}
func (l *LevelVarSLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelError, msg, args...)
}
func (l *LevelVarSLogger) Handler() slog.Handler {
	return l.handler
}
func (l *LevelVarSLogger) Info(msg string, args ...any) {
	l.log(context.Background(), slog.LevelInfo, msg, args...)
}
func (l *LevelVarSLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelInfo, msg, args...)
}
func (l *LevelVarSLogger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	l.log(ctx, level, msg, args...)
}
func (l *LevelVarSLogger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	l.logAttrs(ctx, level, msg, attrs...)
}
func (l *LevelVarSLogger) Warn(msg string, args ...any) {
	l.log(context.Background(), slog.LevelWarn, msg, args...)
}
func (l *LevelVarSLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelWarn, msg, args...)
}

func (l *LevelVarSLogger) With(args ...any) *LevelVarSLogger {
	if len(args) == 0 {
		return l
	}
	c := l.clone()
	c.handler = l.handler.WithAttrs(argsToAttrSlice(args))
	return c
}

func (l *LevelVarSLogger) WithGroup(name string) *LevelVarSLogger {
	if name == "" {
		return l
	}
	c := l.clone()
	c.handler = l.handler.WithGroup(name)
	return c
}
func (l *LevelVarSLogger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if !l.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	if !ignorePC {
		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(3, pcs[:])
		pc = pcs[0]
	}
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}
func (l *LevelVarSLogger) logAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if !l.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	if !ignorePC {
		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(3, pcs[:])
		pc = pcs[0]
	}
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.AddAttrs(attrs...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}
func (l *LevelVarSLogger) clone() *LevelVarSLogger {
	c := *l
	c.level = &slog.LevelVar{}
	c.level.Set(l.level.Level())
	return &c
}
func argsToAttrSlice(args []any) []slog.Attr {
	var (
		attr  slog.Attr
		attrs []slog.Attr
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	return attrs
}
func argsToAttr(args []any) (slog.Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return slog.String("!BADKEY", x), nil
		}
		return slog.Any(x, args[1]), args[2:]

	case slog.Attr:
		return x, args[1:]

	default:
		return slog.Any("!BADKEY", x), args[1:]
	}
}

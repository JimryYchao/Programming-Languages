package examples

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

// ? go test -v -run=^TestLevelVarHandler$
func TestLevelVarHandler(t *testing.T) {
	level := &slog.LevelVar{} // default is info
	lgr := slog.New(NewLevelVarHandler(slog.Default().Handler(), level))
	lgr.Debug("not printed")
	level.Set(slog.LevelDebug - 1)
	lgr.Debug("printed")
}

// ? go test -v -run=^TestLevelVarSLogger$
func TestLevelVarSLogger(t *testing.T) {
	lgr := NewLevelVarSLogger(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(-10000)}), 0)
	lgr.Debug("not printed")
	lgr.Info("printed")
	lgr.Log(context.Background(), -10, "not printed")
	lgr.Log(context.Background(), 10, "test level 10")
	lgr.SetLevel(-100)
	lgr.Debug("printed")

	lgr = lgr.With("author", "Ychao")
	lgr.Debug("printed")

	lgr = lgr.WithGroup("G")
	lgr.Info("printed", "k1", "v1")
}

package logger

import (
	"context"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	"golang.org/x/exp/slog"
)

var _ = pgx.Logger((*PGXAdapter)(nil))

type PGXAdapter struct {
	*slog.Logger
}

func (l *PGXAdapter) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	switch level {
	case pgx.LogLevelTrace:
		l.Logger.DebugCtx(ctx, msg, slog.Any("data", data))
	case pgx.LogLevelDebug:
		l.Logger.DebugCtx(ctx, msg, slog.Any("data", data))
	case pgx.LogLevelInfo:
		l.Logger.InfoCtx(ctx, msg, slog.Any("data", data))
	case pgx.LogLevelWarn:
		l.Logger.WarnCtx(ctx, msg, slog.Any("data", data))
	case pgx.LogLevelError:
		l.Logger.ErrorCtx(ctx, msg, slog.Any("data", data))
	default:
		l.Logger.InfoCtx(ctx, msg, slog.Any("data", data))
	}
}

func NewPGXLogger(l *slog.Logger) *PGXAdapter {
	return &PGXAdapter{
		Logger: l,
	}
}

// CustomJSONHandler tries to get requestID from ctx and use it in structured log.
type CustomJSONHandler struct {
	*slog.JSONHandler
}

func (h *CustomJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	rID := middleware.GetReqID(ctx)
	if rID != "" {
		r.AddAttrs(slog.String("req_id", rID))
	}

	return h.JSONHandler.Handle(ctx, r)
}

func New(level slog.Leveler) *slog.Logger {
	basicJSONHandler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: level,
		},
	)

	logger := slog.New(
		&CustomJSONHandler{
			basicJSONHandler,
		},
	)

	return logger
}

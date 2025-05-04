package mcp

import (
	"context"
	"log/slog"
)

const (
	// Emergency: system is unusable
	LogLevelEmergency slog.Level = 0
	// Alert: action must be taken immediately
	LogLevelAlert slog.Level = 1
	// Critical: critical conditions
	LogLevelCritical slog.Level = 2
	// Error: error conditions
	LogLevelError slog.Level = 3
	// Warning: warning conditions
	LogLevelWarning slog.Level = 4
	// Notice: normal but significant condition
	LogLevelNotice slog.Level = 5
	// Informational: informational messages
	LogLevelInfo slog.Level = 6
	// Debug: debug-level messages
	LogLevelDebug slog.Level = 7
)

func NewLogHandler(t Transport, level slog.Level) slog.Handler {
	return &logHandler{
		t:     t,
		level: level,
	}
}

func convertStrToLevel(level string) slog.Level {
	switch level {
	case "emergency":
		return LogLevelEmergency
	case "alert":
		return LogLevelAlert
	case "critical":
		return LogLevelCritical
	case "error":
		return LogLevelError
	case "warning":
		return LogLevelWarning
	case "notice":
		return LogLevelNotice
	case "info":
		return LogLevelInfo
	case "debug":
		return LogLevelDebug
	default:
		return LogLevelInfo
	}
}

var _ slog.Handler = (*logHandler)(nil)

type logHandler struct {
	t     Transport
	level slog.Level
}

func (h *logHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}
func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.level <= level
}
func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}
func (h *logHandler) WithGroup(name string) slog.Handler {
	return h
}

var _ slog.Handler = (*multiLogHandler)(nil)

type multiLogHandler []slog.Handler

func (h *multiLogHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}
func (h *multiLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}
func (h *multiLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}
func (h *multiLogHandler) WithGroup(name string) slog.Handler {
	return h
}

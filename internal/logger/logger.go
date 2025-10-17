package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Custom log levels
const (
    LevelTrace = slog.Level(-8)
    LevelFatal = slog.Level(12)
)

// Custom log handler
type CustomHandler struct {
    handler slog.Handler
}

func NewCustomHandler(handler slog.Handler) *CustomHandler {
    return &CustomHandler{handler: handler}
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
    // Add custom attributes to all log records
    r.Add("go_version", runtime.Version())
    
    // Add timestamp in ISO format
    r.Add("timestamp", time.Now().Format(time.RFC3339))
    
    return h.handler.Handle(ctx, r)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return NewCustomHandler(h.handler.WithAttrs(attrs))
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
    return NewCustomHandler(h.handler.WithGroup(name))
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
    return h.handler.Enabled(ctx, level)
}

// Logger wrapper with custom methods
type AppLogger struct {
    logger *slog.Logger
}

func NewAppLogger(env string) *AppLogger {
    var handler slog.Handler
    
    // Different handlers for different environments
    switch env {
    case "development":
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level:       LevelTrace,
            AddSource:   true,
        })
    case "production":
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
        })
    default:
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
        })
    }
    
    // Wrap with our custom handler
    customHandler := NewCustomHandler(handler)
    
    logger := slog.New(customHandler)
    
    return &AppLogger{logger: logger}
}

// Custom log methods
func (l *AppLogger) Trace(msg string, args ...any) {
    l.logger.Log(context.Background(), LevelTrace, msg, args...)
}

func (l *AppLogger) Debug(msg string, args ...any) {
    l.logger.Debug(msg, args...)
}

func (l *AppLogger) Info(msg string, args ...any) {
    l.logger.Info(msg, args...)
}

func (l *AppLogger) Warn(msg string, args ...any) {
    l.logger.Warn(msg, args...)
}

func (l *AppLogger) Error(msg string, args ...any) {
    l.logger.Error(msg, args...)
}

func (l *AppLogger) Fatal(msg string, args ...any) {
    l.logger.Log(context.Background(), LevelFatal, msg, args...)
    os.Exit(1)
}

// Context-aware logging
func (l *AppLogger) With(args ...any) *AppLogger {
    return &AppLogger{logger: l.logger.With(args...)}
}

// HTTP request logging
func (l *AppLogger) RequestInfo(method, path, ip string, status int, latency time.Duration, userID string) {
    l.Info("HTTP Request",
        "method", method,
        "path", path,
        "ip", ip,
        "status", status,
        "latency", latency,
        "user_id", userID,
    )
}

// Database query logging
func (l *AppLogger) QueryInfo(query string, duration time.Duration, rowsAffected int64) {
    l.Debug("Database Query",
        "query", query,
        "duration", duration,
        "rows_affected", rowsAffected,
    )
}

// Authentication logging
func (l *AppLogger) AuthInfo(userID, action string) {
    l.Info("Authentication Event",
        "user_id", userID,
        "action", action,
    )
}

// Error with stack trace
func (l *AppLogger) ErrorWithTrace(err error, args ...any) {
    // Get caller information
    pc, file, line, _ := runtime.Caller(1)
    fn := runtime.FuncForPC(pc)
    
    errorArgs := []any{
        "error", err.Error(),
        "file", file,
        "line", line,
        "function", fn.Name(),
    }
    errorArgs = append(errorArgs, args...)
    
    l.Error("Application Error", errorArgs...)
}
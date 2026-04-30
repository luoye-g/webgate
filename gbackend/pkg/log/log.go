package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Config defines the log configuration.
type Config struct {
	// Level: debug | info | warn | error (default: info)
	Level string
	// Format: json | text (default: json)
	Format string
	// Dir: directory for log files. Empty means console only.
	Dir string
	// FileName: base file name, e.g. "app". Final file will be "app-2006-01-02.log"
	FileName string
	// Console: also write to stdout (default: true)
	Console bool
}

// contextKey is the type for context keys used by this package.
type contextKey string

const (
	// TraceIDKey is the context key for trace id.
	TraceIDKey contextKey = "trace_id"
	// UserIDKey is the context key for user id.
	UserIDKey contextKey = "user_id"
)

var (
	defaultLogger *slog.Logger
	rotatingW     *rotatingWriter
	initOnce      sync.Once
)

// Init initializes the global logger. Safe to call once; subsequent calls are no-ops.
func Init(cfg Config) error {
	var err error
	initOnce.Do(func() {
		err = setup(cfg)
	})
	return err
}

func setup(cfg Config) error {
	level := parseLevel(cfg.Level)

	writers := make([]io.Writer, 0, 2)
	if cfg.Console || cfg.Dir == "" {
		writers = append(writers, os.Stdout)
	}

	if cfg.Dir != "" {
		if err := os.MkdirAll(cfg.Dir, 0o755); err != nil {
			return fmt.Errorf("create log dir: %w", err)
		}
		fileName := cfg.FileName
		if fileName == "" {
			fileName = "app"
		}
		rotatingW = newRotatingWriter(cfg.Dir, fileName)
		writers = append(writers, rotatingW)
	}

	var out io.Writer
	if len(writers) == 1 {
		out = writers[0]
	} else {
		out = io.MultiWriter(writers...)
	}

	handlerOpts := &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	}

	var handler slog.Handler
	if strings.EqualFold(cfg.Format, "text") {
		handler = slog.NewTextHandler(out, handlerOpts)
	} else {
		handler = slog.NewJSONHandler(out, handlerOpts)
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
	return nil
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// L returns the global logger. If not initialized, returns a default stdout logger.
func L() *slog.Logger {
	if defaultLogger == nil {
		return slog.Default()
	}
	return defaultLogger
}

// WithContext returns a logger enriched with context fields (trace_id, user_id).
func WithContext(ctx context.Context) *slog.Logger {
	logger := L()
	if ctx == nil {
		return logger
	}
	attrs := make([]any, 0, 4)
	if v := ctx.Value(TraceIDKey); v != nil {
		if s, ok := v.(string); ok && s != "" {
			attrs = append(attrs, slog.String("trace_id", s))
		}
	}
	if v := ctx.Value(UserIDKey); v != nil {
		attrs = append(attrs, slog.Any("user_id", v))
	}
	if len(attrs) == 0 {
		return logger
	}
	return logger.With(attrs...)
}

// Convenience top-level helpers.

func Debug(msg string, args ...any) { L().Debug(msg, args...) }
func Info(msg string, args ...any)  { L().Info(msg, args...) }
func Warn(msg string, args ...any)  { L().Warn(msg, args...) }
func Error(msg string, args ...any) { L().Error(msg, args...) }

// Sync flushes any buffered writer (currently a no-op; reserved for future use).
func Sync() error {
	if rotatingW != nil {
		return rotatingW.Close()
	}
	return nil
}

// rotatingWriter writes to a file whose name is rotated daily.
type rotatingWriter struct {
	mu       sync.Mutex
	dir      string
	baseName string
	curDate  string
	file     *os.File
}

func newRotatingWriter(dir, baseName string) *rotatingWriter {
	return &rotatingWriter{dir: dir, baseName: baseName}
}

func (r *rotatingWriter) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if r.file == nil || r.curDate != today {
		if err := r.rotateLocked(today); err != nil {
			return 0, err
		}
	}
	return r.file.Write(p)
}

func (r *rotatingWriter) rotateLocked(date string) error {
	if r.file != nil {
		_ = r.file.Close()
		r.file = nil
	}
	path := filepath.Join(r.dir, fmt.Sprintf("%s-%s.log", r.baseName, date))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	r.file = f
	r.curDate = date
	return nil
}

// Close closes the underlying file if any.
func (r *rotatingWriter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.file != nil {
		err := r.file.Close()
		r.file = nil
		return err
	}
	return nil
}

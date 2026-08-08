package logging

import (
	"context"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	fileMaxBackups = 10
	fileMaxSize    = 100 // mega-bytes
	fileMaxAge     = 90  // days
)

// Config represents logger settings.
type Config struct {
	FileOut string `split_words:"true"`
	FileErr string `split_words:"true"`
}

// Default returns pre-configured logger.
func Default() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}

// DefaultWithPrefix returns a default instance of logger with prefix.
func DefaultWithPrefix(prefix string) *slog.Logger {
	return WithPrefix(Default(), prefix)
}

// WithPrefix returns a logger with a prefix field.
func WithPrefix(lg *slog.Logger, prefix string) *slog.Logger {
	return lg.With("prefix", prefix)
}

// New returns logger with config.
func New(cfg Config) (*slog.Logger, error) {
	handlers := []slog.Handler{slog.NewTextHandler(os.Stderr, nil)}

	if cfg.FileOut != "" {
		outF, err := openLogfileWithRotator(cfg.FileOut)
		if err != nil {
			return nil, err
		}
		handlers = append(handlers, levelFilterHandler{
			level:   slog.LevelDebug,
			handler: slog.NewTextHandler(outF, nil),
		})
	}

	if cfg.FileErr != "" {
		errF, err := openLogfileWithRotator(cfg.FileErr)
		if err != nil {
			return nil, err
		}
		handlers = append(handlers, levelFilterHandler{
			level:   slog.LevelError,
			handler: slog.NewTextHandler(errF, nil),
		})
	}

	return slog.New(slog.NewMultiHandler(handlers...)), nil
}

// levelFilterHandler only handles records at or above the configured level.
type levelFilterHandler struct {
	level   slog.Leveler
	handler slog.Handler
}

func (h levelFilterHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return l >= h.level.Level() && h.handler.Enabled(ctx, l)
}

func (h levelFilterHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

func (h levelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return levelFilterHandler{h.level, h.handler.WithAttrs(attrs)}
}

func (h levelFilterHandler) WithGroup(name string) slog.Handler {
	return levelFilterHandler{h.level, h.handler.WithGroup(name)}
}

func openLogfileWithRotator(path string) (*lumberjack.Logger, error) {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    fileMaxSize, // megabytes
		MaxBackups: fileMaxBackups,
		MaxAge:     fileMaxAge, // days
		Compress:   true,       // disabled by default
	}, nil
}

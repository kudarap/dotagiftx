package logging

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	fileMaxBackups = 10
	fileMaxSize    = 100 // mega-bytes
	fileMaxAge     = 90  // days
)

const defaultLevel = "info"

// Config represents logger settings.
type Config struct {
	FileOut string `split_words:"true"`
	FileErr string `split_words:"true"`
	Level   string
}

func (c Config) setDefault() Config {
	c.Level = strings.ToLower(strings.TrimSpace(c.Level))
	if c.Level == "" {
		c.Level = defaultLevel
	}
	return c
}

func (c Config) level() (slog.Level, error) {
	v, ok := levels[c.Level]
	if !ok {
		return slog.LevelInfo, fmt.Errorf("un-supported level: %s", c.Level)
	}
	return v, nil
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
	cfg = cfg.setDefault()
	level, err := cfg.level()
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{Level: level}
	handlers := []slog.Handler{slog.NewTextHandler(os.Stderr, opts)}
	if cfg.FileOut != "" {
		outF, err := openLogfileWithRotator(cfg.FileOut)
		if err != nil {
			return nil, err
		}
		handlers = append(handlers, slog.NewTextHandler(outF, opts))
	}
	if cfg.FileErr != "" {
		errF, err := openLogfileWithRotator(cfg.FileErr)
		if err != nil {
			return nil, err
		}
		opts.Level = slog.LevelError
		handlers = append(handlers, slog.NewTextHandler(errF, opts))
	}

	return slog.New(slog.NewMultiHandler(handlers...)), nil
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

var levels = map[string]slog.Level{
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
	"debug": slog.LevelDebug,
}

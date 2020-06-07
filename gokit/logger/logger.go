package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	FileOut string `split_words:"true"`
	FileErr string `split_words:"true"`
}

// Default returns pre-configured logger.
func Default() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return l
}

// New returns logger with config.
func New(cfg Config) (*logrus.Logger, error) {
	l := Default()

	if cfg.FileOut != "" {
		outF, err := openLogfileWithRotate(cfg.FileOut)
		if err != nil {
			return nil, err
		}
		l.AddHook(&Hook{
			Writer: outF,
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
				logrus.DebugLevel,
			},
		})
	}

	if cfg.FileErr != "" {
		errF, err := openLogfileWithRotate(cfg.FileErr)
		if err != nil {
			return nil, err
		}
		l.AddHook(&Hook{
			Writer: errF,
			LogLevels: []logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
			},
		})
	}

	return l, nil
}

func openLogfileWithRotate(path string) (*lumberjack.Logger, error) {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	}, nil
}

func openLogfile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

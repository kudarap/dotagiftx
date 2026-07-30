package logging

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger provides access to logging service.
type Logger interface {
	// Standard Lib logger methods.
	Print(v ...any)
	Println(v ...any)
	Printf(format string, v ...any)

	Panic(v ...any)
	Panicln(v ...any)
	Panicf(format string, v ...any)

	Fatal(v ...any)
	Fatalln(v ...any)
	Fatalf(format string, v ...any)

	// Custom logger methods.
	Info(v ...any)
	Infoln(v ...any)
	Infof(format string, v ...any)

	Error(v ...any)
	Errorln(v ...any)
	Errorf(format string, v ...any)

	Debug(v ...any)
	Debugln(v ...any)
	Debugf(format string, v ...any)

	Warn(v ...any)
	Warnln(v ...any)
	Warnf(format string, v ...any)
}

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
func Default() *logrus.Logger {
	l := logrus.New()
	// l.SetFormatter(&logrus.TextFormatter{
	//	FullTimestamp: true,
	// })
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	l.Formatter = formatter
	return l
}

// DefaultWithPrefix returns a default instance of logger with prefix.
func DefaultWithPrefix(prefix string) *logrus.Entry {
	return WithPrefix(Default(), prefix)
}

func WithPrefix(lg *logrus.Logger, prefix string) *logrus.Entry {
	return lg.WithField("prefix", prefix)
}

// New returns logger with config.
func New(cfg Config) (*logrus.Logger, error) {
	l := Default()

	if cfg.FileOut != "" {
		outF, err := openLogfileWithRotator(cfg.FileOut)
		if err != nil {
			return nil, err
		}
		l.AddHook(&Hook{
			Writer: outF,
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
				logrus.DebugLevel,
				logrus.WarnLevel,
			},
		})
	}

	if cfg.FileErr != "" {
		errF, err := openLogfileWithRotator(cfg.FileErr)
		if err != nil {
			return nil, err
		}
		l.AddHook(&Hook{
			Writer: errF,
			LogLevels: []logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
			},
		})
	}

	return l, nil
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

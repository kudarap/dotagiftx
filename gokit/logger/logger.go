package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger provides access to logging service.
type Logger interface {
	// Standard Lib logger methods.
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicln(v ...interface{})
	Panicf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalln(v ...interface{})
	Fatalf(format string, v ...interface{})

	// Custom logger methods.
	Info(v ...interface{})
	Infoln(v ...interface{})
	Infof(format string, v ...interface{})

	Error(v ...interface{})
	Errorln(v ...interface{})
	Errorf(format string, v ...interface{})

	Debug(v ...interface{})
	Debugln(v ...interface{})
	Debugf(format string, v ...interface{})

	Warn(v ...interface{})
	Warnln(v ...interface{})
	Warnf(format string, v ...interface{})
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
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return l
}

func WithPrefix(prefix string) *logrus.Entry {
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true

	l := Default()
	l.Formatter = formatter
	return l.WithField("prefix", prefix)
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

func openLogfile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

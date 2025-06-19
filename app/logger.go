package app

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

type LogConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	LogFile    string `json:"log_file"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
}

func DefaultConfig() *LogConfig {
	return &LogConfig{
		Level:      "info",
		Format:     "json",
		Output:     "file",
		LogFile:    getDefaultLogPath(),
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 5,
		Compress:   true,
	}
}

func getDefaultLogPath() string {
	var logDir string

	switch runtime.GOOS {
	case "windows":

		if appData := os.Getenv("APPDATA"); appData != "" {
			logDir = filepath.Join(appData, "wallchemy", "logs")
		} else {
			logDir = filepath.Join(os.TempDir(), "wallchemy", "logs")
		}
	case "darwin":

		if home := os.Getenv("HOME"); home != "" {
			logDir = filepath.Join(home, "Library", "Logs", "wallchemy")
		} else {
			logDir = filepath.Join(os.TempDir(), "wallchemy", "logs")
		}
	default:

		if home := os.Getenv("HOME"); home != "" {
			logDir = filepath.Join(home, ".local", "share", "wallchemy", "logs")
		} else {
			logDir = filepath.Join(os.TempDir(), "wallchemy", "logs")
		}
	}

	return filepath.Join(logDir, "app.log")
}

func InitLogger(config *LogConfig) error {
	Logger = logrus.New()

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return err
	}
	Logger.SetLevel(level)

	switch config.Format {
	case "json":
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			PrettyPrint:     false,
		})
	case "text":
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,

			DisableColors: config.Output == "file",
		})
	default:
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}

	var writers []io.Writer

	if config.Output == "console" || config.Output == "both" {
		writers = append(writers, os.Stdout)
	}

	if config.Output == "file" || config.Output == "both" {

		if err := os.MkdirAll(filepath.Dir(config.LogFile), 0755); err != nil {
			return err
		}

		fileWriter := &lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			MaxBackups: config.MaxBackups,
			Compress:   config.Compress,
			LocalTime:  true,
		}
		writers = append(writers, fileWriter)
	}

	if len(writers) > 1 {
		Logger.SetOutput(io.MultiWriter(writers...))
	} else if len(writers) == 1 {
		Logger.SetOutput(writers[0])
	}

	return nil
}

type CallerHook struct {
	logger *logrus.Logger
}

func (hook *CallerHook) Fire(entry *logrus.Entry) error {

	if hook.logger.Level <= logrus.DebugLevel {
		if pc, file, line, ok := runtime.Caller(8); ok {
			funcName := runtime.FuncForPC(pc).Name()
			entry.Data["caller"] = map[string]any{
				"file":     filepath.Base(file),
				"line":     line,
				"function": filepath.Base(funcName),
			}
		}
	}
	return nil
}

func (hook *CallerHook) Levels() []logrus.Level {

	return logrus.AllLevels
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}

func WithField(key string, value any) *logrus.Entry {
	return Logger.WithField(key, value)
}

func WithError(err error) *logrus.Entry {
	return Logger.WithError(err)
}

func LoggerConfigFromEnv() *LogConfig {
	config := DefaultConfig()

	if os.Getenv("DEBUG") != "" || os.Getenv("LOG_TO_STDOUT") != "" {
		config.Output = "console"
		config.Format = "text"
		config.Level = "debug"
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Level = level
	}

	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.Format = format
	}

	if output := os.Getenv("LOG_OUTPUT"); output != "" {
		config.Output = output
	}

	return config
}

func SetupLogger(debug bool) error {
	config := DefaultConfig()

	if debug {
		config.Output = "console"
		config.Format = "text"
		config.Level = "debug"

	} else {
		config.Output = "file"
		config.Format = "json"
		config.Level = "info"
	}

	return InitLogger(config)
}

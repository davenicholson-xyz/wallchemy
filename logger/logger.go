package logger

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	logFile := getLogPath("wallchemy")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.Info("Failed to create log file, using default stderr")
	}

	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.DebugLevel)
}

func getLogPath(appName string) string {
	var logDir string

	switch runtime.GOOS {
	case "linux":
		home, _ := os.UserHomeDir()
		logDir = filepath.Join(home, "local", "share", appName, "logs")
	case "darwin":
		home, _ := os.UserHomeDir()
		logDir = filepath.Join(home, "Library", "Logs", appName)
	case "windows":
		logDir = filepath.Join(os.Getenv("LocalAppData"), appName, "logs")
	default:
		logDir = "./logs"
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	return filepath.Join(logDir, appName+".log")
}

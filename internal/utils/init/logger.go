package inith

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	defaultLogLevel = logrus.DebugLevel
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	log.SetLevel(getLevelFromEnv())
	return log
}

func getLevelFromEnv() logrus.Level {

	levels := map[string]logrus.Level{
		"DEBUG": logrus.DebugLevel,
		"INFO":  logrus.InfoLevel,
		"WARN":  logrus.WarnLevel,
		"ERROR": logrus.ErrorLevel,
		"FATAL": logrus.FatalLevel,
		"PANIC": logrus.PanicLevel,
	}

	currentLevel := os.Getenv("LOG_LEVEL")

	if level, ok := levels[currentLevel]; ok {
		return level
	}
	return defaultLogLevel
}

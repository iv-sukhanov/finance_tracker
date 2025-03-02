package utils

import (
	"github.com/sirupsen/logrus"
)

const (
	defaultLogLevel = logrus.InfoLevel
)

func NewLogger(level string) *logrus.Logger {
	log := (logrus.New())
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	log.SetLevel(GetLevelFromString(level))
	return log
}

func GetLevelFromString(currentLevel string) logrus.Level {

	levels := map[string]logrus.Level{
		"DEBUG": logrus.DebugLevel,
		"INFO":  logrus.InfoLevel,
		"WARN":  logrus.WarnLevel,
		"ERROR": logrus.ErrorLevel,
		"FATAL": logrus.FatalLevel,
		"PANIC": logrus.PanicLevel,
	}

	if level, ok := levels[currentLevel]; ok {
		return level
	}
	return defaultLogLevel
}

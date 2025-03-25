package utils

import (
	"github.com/sirupsen/logrus"
)

const (
	defaultLogLevel = logrus.InfoLevel
)

// NewLogger creates a new logger instance with the specified log level and format.
func NewLogger(level string) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	log.SetLevel(GetLevelFromString(level))
	return log
}

// GetLevelFromString converts a string representation of a log level to a logrus.Level.
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

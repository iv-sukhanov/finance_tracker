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

// GetLevelFromString converts a string representation of a log level
// into the corresponding logrus.Level. Supported log levels are:
// "DEBUG", "INFO", "WARN", "ERROR", "FATAL", and "PANIC".
// If the provided string does not match any of these levels,
// the function returns the defaultLogLevel.
//
// Parameters:
//   - currentLevel: A string representing the desired log level.
//
// Returns:
//   - logrus.Level: The corresponding logrus log level.
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

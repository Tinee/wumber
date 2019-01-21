package logger

import (
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Logger wraps logrus logger.
type Logger struct {
	*log.Logger
}

// NewLogger determines by looking at the environment, which level the logger should log with.
// If the environment isn't set it will default to a Developer logger.
func NewLogger(env string, wr io.Writer) *Logger {
	logger := &Logger{
		Logger: &log.Logger{
			Out: wr,
		},
	}

	switch strings.ToLower(env) {
	case "uat":
		setupUATLogger(logger)
	case "prod":
		setupPRODLogger(logger)
	case "dev":
		setupDEVLogger(logger)
	default:
		setupDEVLogger(logger)
	}

	return logger
}

func setupUATLogger(l *Logger) {
	l.SetFormatter(&log.JSONFormatter{})
	l.SetLevel(log.DebugLevel)
	l.SetReportCaller(true)
}

func setupPRODLogger(l *Logger) {
	l.SetFormatter(&log.JSONFormatter{})
	l.SetLevel(log.ErrorLevel)
}

func setupDEVLogger(l *Logger) {
	l.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
		QuoteEmptyFields: true,
	})
	l.SetLevel(log.TraceLevel)
	l.SetReportCaller(true)
}

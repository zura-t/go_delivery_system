package logger

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Interface -.
type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *log.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level string) *Logger {
	var l log.Level

	switch strings.ToLower(level) {
	case "error":
		l = log.ErrorLevel
	case "warn":
		l = log.WarnLevel
	case "info":
		l = log.InfoLevel
	case "debug":
		l = log.DebugLevel
	default:
		l = log.InfoLevel
	}

	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New()
	logger.SetOutput(file)
	logger.Formatter = &log.JSONFormatter{}
	logger.Hooks = make(log.LevelHooks)
	logger.Level = l

	return &Logger{
		logger: logger,
	}
}

// Debug -.
func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(message, args)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(message, args)
}

// Error -.
func (l *Logger) Error(message interface{}, args ...interface{}) {
	if l.logger.GetLevel() == log.DebugLevel {
		l.logger.Debug(message, args)
	}
	l.logger.Error(message, args)
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)
}

func (l *Logger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info(message)
		l.logger.WithField("message", message)
	} else {
		l.logger.Info(message, args)
		l.logger.WithField("message", message)
	}
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}

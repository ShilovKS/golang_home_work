package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Logger struct {
	logger *log.Logger
	level  Level
}

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var levelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
}

func New(level string) (*Logger, error) {
	logLevel := InfoLevel //nolint:ineffassign
	switch strings.ToLower(level) {
	case "debug":
		logLevel = DebugLevel
	case "info":
		logLevel = InfoLevel
	case "warn", "warning":
		logLevel = WarnLevel
	case "error":
		logLevel = ErrorLevel
	default:
		return nil, fmt.Errorf("unsupported log level: %s", level)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return &Logger{
		logger: logger,
		level:  logLevel,
	}, nil
}

func (l *Logger) log(level Level, msg string) {
	if level < l.level {
		return
	}
	levelName := levelNames[level]
	message := fmt.Sprintf("[%s] %s", levelName, msg)
	l.logger.Println(message)
}

func (l *Logger) Debug(msg string) {
	l.log(DebugLevel, msg)
}

func (l *Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

func (l *Logger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

func (l *Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

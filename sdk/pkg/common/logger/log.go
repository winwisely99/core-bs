package logger

import (
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	*log.Logger
}

func NewLogger(lvl log.Level, fields map[string]interface{}) *Logger {
	l := log.New()
	l.WithFields(fields)
	l.Level = lvl
	return &Logger{l}
}

func (l *Logger) AddFields(fields map[string]interface{}) *Logger {
	l.WithFields(fields)
	return l
}

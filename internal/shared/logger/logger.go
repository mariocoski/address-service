package logger

import (
	"log"
)

type Logger interface {
	Info(message string)
	Error(message string)
}

type standardLogger struct{}

func (l standardLogger) Info(message string) {
	log.Println(message)
}

func (l standardLogger) Error(message string) {
	log.Println(message)
}

func NewLogger() Logger {
	return &standardLogger{}
}

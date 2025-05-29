package logger

import (
	"log"
)

type Logger struct {
	levelLog string
}

func New(level string) *Logger {
	return &Logger{levelLog: level}
}

func (l Logger) Info(msg string) {
	if l.levelLog == "INFO" || l.levelLog == "DEBUG" {
		log.Println("|" + " " + l.levelLog + " " + "|" + " " + msg)
	}
}

func (l Logger) Error(msg string) {
	log.Println("|" + " " + "ERROR" + " " + "|" + " " + msg)
}

func (l Logger) Debug(msg string) {
	if l.levelLog == "DEBUG" {
		log.Println("|" + " " + l.levelLog + " " + "|" + " " + msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.levelLog != "ERROR" {
		log.Println("|" + " " + "WARN" + " " + "|" + " " + msg)
	}
}

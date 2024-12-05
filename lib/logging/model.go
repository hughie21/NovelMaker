package logging

import (
	"sync"
	"time"
)

type Level int

var (
	Logger *Log
	once   *sync.Once
)

const (
	InfoLevel = iota
	WarnLevel
	ErrorLevel
	FatalLevel
	DebugLevel
)

type LogMessage struct {
	Level    Level
	Time     time.Time
	Message  string
	FuncName string
}

type FileLogger struct {
	Filename string
	expired  int
}

type Log struct {
	Level      Level
	Message    []LogMessage
	FileLogger *FileLogger
	FileFlag   bool
}

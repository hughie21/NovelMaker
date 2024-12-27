// Description: Definition of the log file structure
// Author: Hughie21
// Date: 2024-12-20
// license that can be found in the LICENSE file.
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

// Message interface
//
// The Message interface is used to define the methods that must be implemented by the message types.
//
// The String() method returns the message as a string.
//
// The getLevel() method returns the log level of the message.
type Message interface {
	String() string
	getLevel() Level
}

type LogMessage struct {
	Level    Level
	Time     time.Time
	Message  string
	FuncName string
}

type TraceMessage struct {
	Time    time.Time
	Message string
	Stack   string
}

type FileLogger struct {
	Filename string
	expired  int
}

type Log struct {
	Level      Level
	Message    []Message
	FileLogger *FileLogger
	FileFlag   bool
}

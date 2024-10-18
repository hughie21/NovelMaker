package logging

import (
	sys "NovelMaker/sys"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Level map[string]int

var (
	Logger     *Log
	once       *sync.Once
	InfoLevel  = Level{"weight": 1}
	WarnLevel  = Level{"weight": 2}
	ErrorLevel = Level{"weight": 3}
	DebugLevel = Level{"weight": 5}
	FatalLevel = Level{"weight": 4}
)

type LogMessage struct {
	Level    Level
	Time     time.Time
	Message  string
	FuncName string
}

type FileLogger struct {
	Filename string
}

type Log struct {
	Level      Level
	Message    []LogMessage
	FileLogger *FileLogger
	FileFlag   bool
}

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func (message *LogMessage) toString() string {
	timeFormat := message.Time.Format("2006-01-02 15:04:05")
	var str_level string
	switch message.Level["weight"] {
	case 1:
		str_level = "INFO"
	case 2:
		str_level = "WARN"
	case 3:
		str_level = "ERROR"
	case 5:
		str_level = "DEBUG"
	case 4:
		str_level = "FATAL"
	}
	return fmt.Sprintf("%s [%s] %s -> %s \n", timeFormat, str_level, message.FuncName, message.Message)
}

func (fl *FileLogger) Print(message string) error {
	f, err := os.OpenFile(fl.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		sys.ShowMessage("Warning", "Unable to write the log file: "+err.Error(), "warning")
		return err
	}
	f.WriteString(message)
	defer f.Close()
	return nil
}

func NewLog(level Level, file bool) *Log {
	if once == nil {
		once = &sync.Once{}
	}
	once.Do(func() {
		Logger = &Log{
			Level:      level,
			FileLogger: &FileLogger{},
			FileFlag:   file,
		}
	})
	return Logger
}

func (l *Log) SetLevel(level Level) {
	l.Level = level
}

func (l *Log) SetFileLogger(filename string) {
	print(filename)
	l.FileLogger.Filename = filename
}

func (l *Log) AddLogMessage(logMessage LogMessage) {
	l.Message = append(l.Message, logMessage)
}

func (l *Log) Info(message string, funcName string) {
	if message != "" {
		l.AddLogMessage(LogMessage{
			Level:    InfoLevel,
			Time:     time.Now(),
			Message:  message,
			FuncName: funcName,
		})
	}
}

func (l *Log) Warning(message string, funcName string) {
	if message != "" {
		l.AddLogMessage(LogMessage{
			Level:    WarnLevel,
			Time:     time.Now(),
			Message:  message,
			FuncName: funcName,
		})
	}
}

func (l *Log) Debug(message string, funcName string) {
	if message != "" {
		l.AddLogMessage(LogMessage{
			Level:    DebugLevel,
			Time:     time.Now(),
			Message:  message,
			FuncName: funcName,
		})
	}
}

func (l *Log) Fatal(message string, funcName string) {
	if message != "" {
		l.AddLogMessage(LogMessage{
			Level:    FatalLevel,
			Time:     time.Now(),
			Message:  message,
			FuncName: funcName,
		})
	}
}

func (l *Log) Error(message string, funcName string) {
	if message != "" {
		l.AddLogMessage(LogMessage{
			Level:    ErrorLevel,
			Time:     time.Now(),
			Message:  message,
			FuncName: funcName,
		})
	}
}

func (l *Log) LogOutPut(rootpath string) error {
	if !l.FileFlag {
		return nil
	}
	todaystr := time.Now().Format("2006-01-02")
	todayint, _ := time.ParseInLocation("2006-01-02", todaystr, time.Local)
	logFileName := strconv.FormatInt(todayint.Unix(), 10) + ".log"
	l.SetFileLogger(filepath.Join(rootpath, "log", logFileName))
	var err error
	for _, message := range l.Message {
		fmt.Println(message.Message)
		if message.Level["weight"] <= l.Level["weight"] {
			err = l.FileLogger.Print(message.toString())
		}
	}
	return err
}

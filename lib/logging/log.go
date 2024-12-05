package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hughie21/NovelMaker/lib/utils"
)

func RunFuncName() string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d | %s", file, line, f.Name())
}

func (message *LogMessage) String() string {
	timeFormat := message.Time.Format("2006-01-02 15:04:05")
	var str_level string
	switch message.Level {
	case 0:
		str_level = "INFO"
	case 1:
		str_level = "WARN"
	case 2:
		str_level = "ERROR"
	case 3:
		str_level = "DEBUG"
	case 4:
		str_level = "FATAL"
	}
	return fmt.Sprintf("%s [%s] %s -> %s \n", timeFormat, str_level, message.FuncName, message.Message)
}

func (fl *FileLogger) Print(message string) error {
	f, err := os.OpenFile(fl.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	f.WriteString(message)
	defer f.Close()
	return nil
}

func NewLog(level Level, file bool, expired int) *Log {
	if once == nil {
		once = &sync.Once{}
	}
	once.Do(func() {
		Logger = &Log{
			Level: level,
			FileLogger: &FileLogger{
				expired: expired,
			},
			FileFlag: file,
		}
	})
	return Logger
}

func (l *Log) SetLevel(level Level) {
	l.Level = level
}

func (l *Log) Expired() {

}

func (l *Log) SetFileLogger(filename string) {
	if !utils.PathExists(filename) {
		fs, err := os.Create(filename)
		if err != nil {
			utils.ShowMessage("Error", "Failed to create log file: "+err.Error(), "error")
		}
		defer fs.Close()
		systemInfo := NewSystem()
		fs.WriteString(systemInfo.String() + "\n")
	}
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

func (l *Log) expired(rootpath string) error {
	expireTime := l.FileLogger.expired // days
	expireDuration := time.Duration(expireTime) * 24 * time.Hour
	now := time.Now()
	err := filepath.Walk(filepath.Join(rootpath, "log"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".log" {
			filename := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			timeInt, err := strconv.ParseInt(filename, 10, 64)
			fileTime := time.Unix(timeInt, 0)
			if err != nil {
				return nil
			}
			if now.Sub(fileTime) > expireDuration {
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func (l *Log) LogOutPut(rootpath string) error {
	err := l.expired(rootpath)
	if err != nil {
		return err
	}
	if !l.FileFlag {
		return nil
	}
	todaystr := time.Now().Format("2006-01-02")
	todayint, _ := time.ParseInLocation("2006-01-02", todaystr, time.Local)
	logFileName := strconv.FormatInt(todayint.Unix(), 10) + ".log"
	l.SetFileLogger(filepath.Join(rootpath, "log", logFileName))
	for _, message := range l.Message {
		fmt.Println(message.Message)
		if message.Level <= l.Level {
			err = l.FileLogger.Print(message.String())
			if err != nil {
				return err
			}
		}
	}
	return err
}

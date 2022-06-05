package mlogger

import (
	"encoding/json"
	"fmt"
	mlogger "github.com/masuldev/mlogger/rotate"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

type Logger struct {
	worker *log.Logger
}

type LogInfo struct {
	Timestamp string
	Level     string
	Caller    string
	Message   string
}

func NewDefaultLogger() (*Logger, error) {
	defaultRotate, err := mlogger.NewDefaultRotate()
	if err != nil {
		return nil, errors.Wrap(err, "Err Create Default RotateLog")
	}

	multiWriter := io.MultiWriter(defaultRotate, os.Stdout)

	return &Logger{worker: log.New(multiWriter, "", 0)}, nil

}

func NewLogger(writer io.Writer) (*Logger, error) {
	var multiWriter io.Writer
	if writer != nil {
		multiWriter = io.MultiWriter(writer, os.Stdout)
	} else {
		multiWriter = io.MultiWriter(nil, os.Stdout)
	}

	return &Logger{worker: log.New(multiWriter, "", 0)}, nil
}

func makeTimestamp() string {
	timeFormat := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Seoul")
	return time.Now().In(loc).Format(timeFormat)
}

func (l *Logger) logging(level int, message string) error {
	_, filename, line, _ := runtime.Caller(2)
	filename = path.Base(filename)

	info := &LogInfo{
		Timestamp: makeTimestamp(),
		Level:     logLevelString(level),
		Caller:    fmt.Sprintf("%s:%v", filename, line),
		Message:   message,
	}

	bytes, _ := json.Marshal(info)

	return l.worker.Output(3, string(bytes))
}

func (l *Logger) Debug(message string) {
	l.logging(0, message)
}

func (l *Logger) Info(message string) {
	l.logging(1, message)
}

func (l *Logger) Warning(message string) {
	l.logging(2, message)
}

func (l *Logger) Error(message string) {
	l.logging(3, message)
}

func (l *Logger) Critical(message string) {
	l.logging(4, message)
}

func logLevelString(level int) string {
	logLevels := [...]string{
		"DEBUG",
		"INFO",
		"WARNING",
		"ERROR",
		"CRITICAL",
	}

	return logLevels[level]
}

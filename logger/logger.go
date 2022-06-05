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
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	CriticalLevel
)

type Logger struct {
	worker *log.Logger
}

type LogInfo struct {
	Level   string
	Caller  string
	Message string
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

func (l *Logger) logging(level int, message string) error {
	_, filename, line, _ := runtime.Caller(2)
	filename = path.Base(filename)

	info := &LogInfo{
		Level:   logLevelString(level),
		Caller:  fmt.Sprintf("%s:%v", filename, line),
		Message: message,
	}

	bytes, _ := json.Marshal(info)

	return l.worker.Output(3, string(bytes))
}

func (l *Logger) Debug(message string) {
	l.logging(0, message)
}

func logLevelString(level int) string {
	logLevels := [...]string{
		"CRITICAL",
		"ERROR",
		"WARNING",
		"INFO",
		"DEBUG",
	}

	return logLevels[level]
}

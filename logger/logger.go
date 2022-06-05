package mlogger

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/lestrrat-go/strftime/internal/errors"
)

type LogLevel int

const (
	CriticalLevel LogLevel = iota + 1
	ErrorLevel
	WarningLevel
	InfoLevel
	DebugLevel
)

type Logger struct {
	worker *log.Logger
}

type LogInfo struct {
	Level     string
	Timestamp time.Time
	Caller    string
	Message   string
}

func NewDefaultLogger() (*Logger, error) {
	defaultRotate, err := NewDefaultRotate()
	if err != nil {
		return nil, errors.Errorf("Create Default RotateLog %s", err)
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

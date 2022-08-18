package mlogger

import (
	"github.com/goccy/go-json"
	mlogger "github.com/masuldev/mlogger/rotate"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	worker *log.Logger
	depth  int
}

type LogInfo struct {
	Timestamp string
	Level     string
	Caller    string
	Message   string
}

func NewDefaultLogger(depth int) (*Logger, error) {
	defaultRotate, err := mlogger.NewDefaultRotate()
	if err != nil {
		return nil, errors.Wrap(err, "Err Create Default RotateLog")
	}

	return NewLogger(defaultRotate, depth)
}

func NewLogger(writer io.Writer, depth int) (*Logger, error) {
	var multiWriter io.Writer
	if writer != nil {
		multiWriter = io.MultiWriter(writer, os.Stdout)
	} else {
		multiWriter = io.Writer(os.Stdout)
	}

	return &Logger{worker: log.New(multiWriter, "", 0), depth: depth}, nil
}

func makeTimestamp() string {
	timeFormat := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Seoul")
	return time.Now().In(loc).Format(timeFormat)
}

func (l *Logger) logging(level int, message string) error {
	//file, line, _ := getActualStack()

	info := &LogInfo{
		Timestamp: makeTimestamp(),
		Level:     logLevelString(level),
		//Caller:    fmt.Sprintf("%s:%v", file, line),
		Message: message,
	}

	bytes, _ := json.Marshal(info)

	return l.worker.Output(l.depth, string(bytes))
}

func getActualStack() (file string, line int, ok bool) {
	cpc, _, _, ok := runtime.Caller(2)
	if !ok {
		return
	}

	callerFunPtr := runtime.FuncForPC(cpc)
	if callerFunPtr == nil {
		ok = false
		return
	}

	var pc uintptr
	for callLevel := 3; callLevel < 5; callLevel++ {
		pc, file, line, ok = runtime.Caller(callLevel)
		file = path.Base(file)
		if !ok {
			return
		}
		funcPtr := runtime.FuncForPC(pc)
		if funcPtr == nil {
			ok = false
			return
		}
		if getFuncNameWithoutPackage(funcPtr.Name()) != getFuncNameWithoutPackage(callerFunPtr.Name()) {
			return
		}
	}
	ok = false
	return
}

func getFuncNameWithoutPackage(name string) string {
	pos := strings.LastIndex(name, ".")
	if pos >= 0 {
		name = name[pos+1:]
	}
	return name
}

func messageMarshaling(message interface{}) string {
	marshalMessage, _ := json.Marshal(message)
	return string(marshalMessage)
}

func (l *Logger) Debug(message interface{}) {
	l.logging(0, messageMarshaling(message))
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

func (l *Logger) Panic(message string) {
	l.logging(5, message)
	os.Exit(1)
}

func logLevelString(level int) string {
	logLevels := [...]string{
		"DEBUG",
		"INFO",
		"WARNING",
		"ERROR",
		"CRITICAL",
		"Panic",
	}

	return logLevels[level]
}

package mlogger

import (
	"fmt"
	"github.com/goccy/go-json"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	mu        sync.RWMutex
	color     bool
	out       FdWriter
	debug     bool
	timestamp bool
}

type LogInfo struct {
	Timestamp string
	Level     string
	Caller    string
	Message   string
}

type FdWriter interface {
	io.Writer
	Fd() uintptr
}

//func NewDefaultLogger(depth int) (*Logger, error) {
//	defaultRotate, err := mlogger.NewDefaultRotate()
//	if err != nil {
//		return nil, errors.Wrap(err, "Err Create Default RotateLog")
//	}
//
//	return NewLogger(defaultRotate)
//}

func NewLogger(out FdWriter) *Logger {
	return &Logger{
		color:     terminal.IsTerminal(int(out.Fd())),
		out:       out,
		timestamp: true,
	}
}

func makeTimestamp() string {
	timeFormat := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Seoul")
	return time.Now().In(loc).Format(timeFormat)
}

func (l *Logger) logging(depth int, data string) error {
	file, line, _ := getActualStack(depth)

	info := &LogInfo{
		Timestamp: makeTimestamp(),
		Level:     logLevelString(depth),
		Caller:    fmt.Sprintf("%s:%v", file, line),
		Message:   data,
	}

	bytes, _ := json.Marshal(info)

	_, err := l.out.Write(bytes)
	return err
}

func getActualStack(level int) (file string, line int, ok bool) {
	cpc, _, _, ok := runtime.Caller(level)
	if !ok {
		return
	}

	callerFunPtr := runtime.FuncForPC(cpc)
	if callerFunPtr == nil {
		ok = false
		return
	}

	var pc uintptr
	for callLevel := level + 1; callLevel < 5; callLevel++ {
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
	var data map[string]interface{}

	marshalMessage, _ := json.Marshal(message)
	json.Unmarshal(marshalMessage, &data)

	for key, value := range data {
		fmt.Println(key, value)
	}

	return string(marshalMessage)
}

func (l *Logger) Debug(v ...interface{}) {
	l.logging(1, fmt.Sprintln(v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.logging(1, fmt.Sprintln(v...))
}

func (l *Logger) Warning(v ...interface{}) {
	l.logging(1, fmt.Sprintln(v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.logging(1, fmt.Sprintln(v...))
}

func (l *Logger) Critical(v ...interface{}) {
	l.logging(1, fmt.Sprintln(v...))
}

func (l *Logger) Panic(v ...interface{}) {
	l.logging(1, fmt.Sprintln(v...))
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

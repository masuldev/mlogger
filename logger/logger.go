package mlogger

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/masuldev/mlogger/internal/buffer"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Logger struct {
	mu        sync.RWMutex
	color     bool
	out       FdWriter
	debug     bool
	timestamp bool
	buf       buffer.Buffer
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

func (l *Logger) Output(depth int, level int, data string) error {
	_, file, line, _ := runtime.Caller(depth + 1)
	file = filepath.Base(file)

	info := &LogInfo{
		Timestamp: makeTimestamp(),
		Level:     logLevelString(level),
		Caller:    fmt.Sprintf("%s:%v", file, line),
		Message:   data,
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.buf.Flush()

	d, _ := json.Marshal(info)
	l.buf.Append(d)

	if len(d) == 0 || d[len(d)-1] != '\n' {
		l.buf.AppendByte('\n')
	}

	_, err := l.out.Write(l.buf)
	return err
}

//func getActualStack(depth int) (file string, line int, ok bool) {
//	cpc, _, _, ok := runtime.Caller(depth)
//	if !ok {
//		return
//	}
//
//	callerFunPtr := runtime.FuncForPC(cpc)
//	if callerFunPtr == nil {
//		ok = false
//		return
//	}
//
//	var pc uintptr
//	for callLevel := depth + 1; callLevel < 5; callLevel++ {
//		pc, file, line, ok = runtime.Caller(callLevel)
//		file = path.Base(file)
//		if !ok {
//			return
//		}
//		funcPtr := runtime.FuncForPC(pc)
//		if funcPtr == nil {
//			ok = false
//			return
//		}
//		if getFuncNameWithoutPackage(funcPtr.Name()) != getFuncNameWithoutPackage(callerFunPtr.Name()) {
//			return
//		}
//	}
//	ok = false
//	return
//}
//
//func getFuncNameWithoutPackage(name string) string {
//	pos := strings.LastIndex(name, ".")
//	if pos >= 0 {
//		name = name[pos+1:]
//	}
//	return name
//}

func (l *Logger) Debug(v ...interface{}) {
	l.Output(1, 0, fmt.Sprint(v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.Output(1, 1, fmt.Sprint(v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.Output(1, 2, fmt.Sprint(v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.Output(1, 3, fmt.Sprint(v...))
}

func (l *Logger) Critical(v ...interface{}) {
	l.Output(1, 4, fmt.Sprint(v...))
}

func (l *Logger) Panic(v ...interface{}) {
	l.Output(1, 5, fmt.Sprint(v...))
	os.Exit(1)
}

func logLevelString(level int) string {
	logLevels := [...]string{
		"DEBUG",
		"INFO",
		"WARNING",
		"ERROR",
		"CRITICAL",
		"PANIC",
	}

	return logLevels[level]
}

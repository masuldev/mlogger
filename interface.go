package mlogger

import (
	"github.com/lestrrat-go/strftime"
	"os"
	"sync"
	"time"
)

type Handler interface {
	Handle(Event)
}

type HandlerFunc func(Event)

type Event interface {
	Type() EventType
}

type EventType int

const (
	InvalidEventType EventType = iota
	FileRotatedEventType
)

type FileRotatedEvent struct {
	prev    string
	current string
}

type Clock interface {
	Now() time.Time
}

type RotateLog struct {
	clock         Clock
	curFn         string
	curBaseFn     string
	globalPattern string
	generation    int
	linkName      string
	maxAge        time.Duration
	mutex         sync.RWMutex
	eventHandler  Handler
	outFh         *os.File
	pattern       *strftime.Strftime
	rotationTime  time.Duration
	rotationSize  int64
	rotationCount uint
	forceNewFile  bool
}

type clockFn func() time.Time

var KST = clockFn(func() time.Time {
	loc, _ := time.LoadLocation("Asia/Seoul")
	return time.Now().In(loc)
})

var Local = clockFn(time.Now)

type Option interface {
	Name() string
	Value() interface{}
}

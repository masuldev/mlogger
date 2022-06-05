package mlogger

import (
	"github.com/masuldev/mlogger/internal/option"
	"time"
)

const (
	optionClock         = "clock"
	optionHandler       = "handler"
	optionLinkName      = "link-name"
	optionMaxAge        = "max-age"
	optionRotationTime  = "rotation-time"
	optionRotationSize  = "rotation-size"
	optionRotationCount = "rotation-count"
	optionForceNewFile  = "force-new-file"
)

func WithClock(c Clock) Option {
	return option.New(optionClock, c)
}

func WithKSTLocation() Option {
	return option.New(optionClock, KST)
}

func WithLinkName(s string) Option {
	return option.New(optionLinkName, s)
}

func WithMaxAge(d time.Duration) Option {
	return option.New(optionMaxAge, d)
}

func WithRotationTime(d time.Duration) Option {
	return option.New(optionRotationTime, d)
}

func WithRotationSize(s int64) Option {
	return option.New(optionRotationSize, s)
}

func WithRotationCount(n uint) Option {
	return option.New(optionRotationCount, n)
}

func WithHandler(h Handler) Option {
	return option.New(optionHandler, h)
}

func ForceNewFile() Option {
	return option.New(optionForceNewFile, true)
}
package mlogger

import (
	"testing"
)

func TestLoggerOutput(t *testing.T) {
	logger, _ := NewLogger(nil, 2)
	logger.Info("this is Info")
	logger.Debug("this is Debug")
	logger.Error("this is Error")
	logger.Warning("this is Warning")
	logger.Critical("this is Critical")
}

package mlogger

import (
	"testing"
)

type TestStruct struct {
	Name    string
	Message string
}

func TestLoggerOutput(t *testing.T) {
	logger, _ := NewLogger(nil, 2)
	logger.Info("this is Info")
	logger.Debug(TestStruct{Name: "Hello", Message: "모두들 안녕"})
	logger.Debug("this is Debug")
	logger.Error("this is Error")
	logger.Warning("this is Warning")
	logger.Critical("this is Critical")
}

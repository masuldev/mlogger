package mlogger

import (
	"os"
	"testing"
)

type TestStruct struct {
	Name    string
	Message string
}

func TestLoggerOutput(t *testing.T) {
	logger := NewLogger(os.Stdout)
	logger.Debug(TestStruct{Name: "Hello", Message: "모두들 안녕"})
	logger.Info("String")
	logger.Warn("String")
	logger.Error("String")
	logger.Panic("String")
}

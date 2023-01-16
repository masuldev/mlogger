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
	logger := NewLogger(os.Stderr)
	//logger.Debug(TestStruct{Name: "Hello", Message: "모두들 안녕"})
	logger.Debug("String")
	logger.Debug("String")
	logger.Debug("String")
	logger.Debug("String")
	//logger.Debug([]string{"Hello", "Array"})
	//logger.Debug(11223123)
}

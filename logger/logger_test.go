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
	//logger.Debug(TestStruct{Name: "Hello", Message: "모두들 안녕"})
	logger.Debug("String")
	logger.Debug("String")
	logger.Debug("String")
	logger.Debug("String")
	//logger.Debug([]string{"Hello", "Array"})
	//logger.Debug(11223123)
}

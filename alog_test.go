package alog

import "testing"

func TestLogger(t *testing.T) {
	var appender ConsoleAppender
	appender.MessageQueue = make(chan string)
	logger := NewLogger("testlogger", []Appender{appender})
	if logger == nil {
		t.Fail()
	}
}

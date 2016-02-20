package alog

import (
	"strings"
	"testing"
	"time"
)

// A test appender that adds messages to a list
type TestAppender struct {
	MessageQueue chan string
	quit         chan int
	logs         []string
}

func (app *TestAppender) StartListening(l *Logger) {
	for {
		select {
		case msg := <-app.MessageQueue:
			app.logs = append(app.logs, msg)
		case <-app.quit:
			l.wg.Done()
			return
		}
	}
}

func (app *TestAppender) Append(msg string) {
	app.MessageQueue <- msg
}

// Finalize triggers the end of the logging loop
func (app *TestAppender) Finalize() {
	app.quit <- 1
}
func TestLogger(t *testing.T) {
	var appender TestAppender
	appender.MessageQueue = make(chan string)
	appender.logs = make([]string, 0)
	appender.quit = make(chan int)
	logger := NewLogger("testlogger", []Appender{&appender})
	if logger == nil {
		t.Fail()
	}
	if len(logger.appenders) != 1 {
		t.Fail()
	}
	logger.Log("Test message", INFO)
	// Wait a reasonable time for the message to be logged
	time.Sleep(100 * time.Millisecond)
	if len(appender.logs) != 1 {
		t.Fail()
	}
	if strings.Compare(appender.logs[0], "Test message") != 0 {
		t.Fail()
	}
	logger.Finalize()
}

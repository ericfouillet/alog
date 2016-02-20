/*
Package alog implements an asynchronous logger.

Basic appenders are included, but the package is designed
to be extended with extra appenders.
*/
package alog

import "sync"

type Logger struct {
	name      string
	appenders []Appender
	listener  chan string
	wg        sync.WaitGroup
	quit      chan int
}

func NewLogger(name string, appenders []Appender) *Logger {
	logger := new(Logger)
	logger.name = name
	logger.appenders = make([]Appender, 0)
	logger.listener = make(chan string) // TODO the channel needs to be closed
	logger.quit = make(chan int)
	for _, appender := range appenders {
		logger.AddAppender(appender)
	}
	logger.wg.Add(1)
	go logger.DispatchMessages()
	return logger
}

type Level int

const (
	INFO  Level = iota
	DEBUG Level = iota
	WARN  Level = iota
	ERROR Level = iota
)

func (logger *Logger) DispatchMessages() {
	for {
		select {
		case msg := <-logger.listener:
			for _, appender := range logger.appenders {
				appender.Append(msg)
			}
		case <-logger.quit:
			logger.wg.Done()
			return
		}
	}
}

func (l *Logger) Finalize() {
	l.quit <- 1
	for _, app := range l.appenders {
		app.Finalize()
	}
	l.wg.Wait()
}

func (l *Logger) Log(msg string, level Level) {
	l.listener <- msg
}

func (l *Logger) AddAppender(app Appender) {
	l.appenders = append(l.appenders, app)
	l.wg.Add(1)
	go app.StartListening(l)
}

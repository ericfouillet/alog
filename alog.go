/*
Package alog implements an asynchronous logger.

Basic appenders are included, but the package is designed
to be extended with extra appenders.
*/
package alog

import (
	"log"
	"sync"
)

////////////// Logger

// A Logger. Is defined by a name and can be linked with multiple appenders
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
	//for _, appender := range logger.appenders {
	//	logger.wg.Add(1)
	//	go appender.StartListening(logger)
	//}
	logger.wg.Add(1)
	go logger.DispatchMessages()
	return logger
}

// Logging level
type Level int

const (
	INFO  Level = iota
	DEBUG Level = iota
	WARN  Level = iota
	ERROR Level = iota
)

func (logger *Logger) DispatchMessages() {
	//log.Println("Dispatching message: getting it from logger.listener")
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
		//logger.wg.Done()
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
	//log.Println("Putting new message in the queue. " + msg)
	l.listener <- msg
}

func (l *Logger) AddAppender(app Appender) {
	l.appenders = append(l.appenders, app)
	l.wg.Add(1)
	go app.StartListening(l)
}

///////////// Appender

// An appender: defines a way to print logs.
// The Append method is asynchronous: the message is simply added to a queue,
// ready to be logged by the logging routine.
type Appender interface {
	Append(msg string)
	StartListening(l *Logger)
	Finalize()
}

// A console appender.
type ConsoleAppender struct {
	MessageQueue chan string
	quit         chan int
}

// Listen forever for log events
func (app ConsoleAppender) StartListening(l *Logger) {
	for {
		select {
		case msg := <-app.MessageQueue:
			//log.Println("Picking a message from the queue")
			log.Println(msg)
		case <-app.quit:
			l.wg.Done()
			return
		}
	}
}

func (app ConsoleAppender) Append(msg string) {
	//log.Println("Putting a message in the appender queue " + msg)
	app.MessageQueue <- msg
}

// Finalize triggers the end of the logging loop
func (app ConsoleAppender) Finalize() {
	app.quit <- 1
}

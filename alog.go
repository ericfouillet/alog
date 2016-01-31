package alog

import (
	"log"
	"sync"
)

////////////// Logger

var wg sync.WaitGroup

// A Logger. Is defined by a name and can be linked with multiple appenders
type Logger struct {
	name      string
	appenders []*Appender
	listener  chan string
	wg        sync.WaitGroup
}

func (logger *Logger) Finalize() {
	logger.wg.Wait()
}

func NewLogger(name string) *Logger {
	logger := &Logger{name, nil, make(chan string), wg}
	logger.AddAppender(&ConsoleAppender{make(chan string), logger.wg})
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
	msg := <-logger.listener
	for _, appender := range logger.appenders {
		logger.wg.Add(1)
		(*appender).Append(msg)
	}
}

func (l *Logger) Println(msg string, level Level) {
	//log.Println("Putting new message in the queue. " + msg)
	go l.DispatchMessages()
	l.listener <- msg
	l.wg.Wait()
}

func (l *Logger) AddAppender(app Appender) {
	go app.StartListening()
	l.appenders = append(l.appenders, &app)
	// Start a goroutine for this appender to print logs
}

///////////// Appender

// An appender: defines a way to print logs.
// The Append method is asynchronous: the message is simply added to a queue,
// ready to be logged by the logging routine.
type Appender interface {
	Append(msg string)
	StartListening()
}

// A console appender.
type ConsoleAppender struct {
	messageQueue chan string
	wg           sync.WaitGroup
}

// Listen forever for log events
func (app *ConsoleAppender) StartListening() {
	for {
		//log.Println("Picking a message from the queue")
		msg := <-app.messageQueue
		log.Println(msg)
		app.wg.Done()
	}
}

func (app *ConsoleAppender) Append(msg string) {
	//log.Println("Putting a message in the appender queue " + msg)
	app.messageQueue <- msg
}

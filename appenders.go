package alog

import "log"

// An Appender: The Append method is asynchronous: the message is simply added to a queue,
// ready to be logged by the logging routine.
type Appender interface {
	Append(msg string)
	StartListening(l *Logger)
	Finalize()
}

type ConsoleAppender struct {
	MessageQueue chan string
	quit         chan int
}

func (app *ConsoleAppender) StartListening(l *Logger) {
	for {
		select {
		case msg := <-app.MessageQueue:
			log.Println(msg)
		case <-app.quit:
			l.wg.Done()
			return
		}
	}
}

func (app *ConsoleAppender) Append(msg string) {
	app.MessageQueue <- msg
}

// Finalize triggers the end of the logging loop
func (app *ConsoleAppender) Finalize() {
	app.quit <- 1
}

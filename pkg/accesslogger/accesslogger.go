package accesslogger

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	BatchSize = 2
	BatchTimeout = "2s"
)
var (
	logger *accessLogger
	batchTimeoutDuration, _ = time.ParseDuration(BatchTimeout)
)

func Start() {
	f, err := os.Create("log")
	if err != nil {
		panic(err)
	}

	logger = &accessLogger{
		events:      make(chan *http.Request, 100),
		stoppedChan: make(chan bool),
		file:        f,
	}

	go func() {
		logger.process()
	}()
}

func Stop() {
	logger.Stop()
}

func Push(req *http.Request) {
	logger.push(req)
}

type accessLogger struct {
	events chan *http.Request
	stoppedChan chan bool
	file *os.File
}

func (l *accessLogger) Stop() {
	close(l.events)
	<- l.stoppedChan
	l.file.Close()
}

func (l *accessLogger) push(req *http.Request) {
	l.events <- req
}

func (l *accessLogger) process() {
	var requests []*http.Request

	for {
		select {
			case request, ok := <- l.events:
				if !ok {
					l.stoppedChan <- true
					return
				}

				requests = append(requests, request)
				if len(requests) >= BatchSize {
					l.processBatch(requests)
					requests = nil
				}
			case <- time.After(batchTimeoutDuration):
				if len(requests) > 0 {
					l.processBatch(requests)
					requests = nil
				}
		}
	}
}

func (l *accessLogger) processBatch(requests []*http.Request) {
	fmt.Printf(">> Processing batch of %d\n", len(requests))
	for _, request := range requests {
		l.file.WriteString(fmt.Sprintf("%s\n", request.URL.String()))
	}

	l.file.Sync()
}

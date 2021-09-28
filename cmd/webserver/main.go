package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webserver/pkg/accesslogger"
	"webserver/pkg/server"
)

func main() {
	fmt.Println("Hello world")

	slice := []string{"a", "b"}

	for _, el := range slice {
		fmt.Println(el)
	}
	//var srv *server.Server
	//(*srv).ListenAndServe()

	accesslogger.Start()
	srv := server.New()
	srv.ListenAndServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<- stop

	fmt.Println("Grace started")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	srv.Shutdown(ctx)
	cancel()

	accesslogger.Stop()

	fmt.Println("Grace complete")
}

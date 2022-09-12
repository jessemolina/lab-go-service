package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"


    "go.uber.org/automaxprocs/maxprocs"
)

var build = "develop"

func main() {

	// ================================================================
	// GOROUTINES
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("maxprocs: %w", err)
		os.Exit(1)
	}
	g := runtime.GOMAXPROCS(0)

	// ================================================================
	// LOGGING

	log.Printf("starting service build[%s] CPU[%d]", build, g)
	defer log.Println("service ended")

	// ================================================================
	// SHUTDOWN

	// make a channel with 1 buffer for an os.Signal
	// block on the channel until it receives either SIGINT or SIGTERM
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	fmt.Println("stopping service")

}

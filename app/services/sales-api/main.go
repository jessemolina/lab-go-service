package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var build = "develop"

func main() {

	// ================================================================
	// LOGGING

	log.Println("starting service", build)
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

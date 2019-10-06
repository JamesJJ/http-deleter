package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func gracefulStop() {

	// Handle ^C and SIGTERM gracefully
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Fprintf(os.Stderr, "Killed: %+v", sig)
		os.Exit(0)
	}()
}

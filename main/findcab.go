package main

import (
	"flag"
	"github.com/gyokuro/findcab"
	"io"
)

func main() {

	flag.Parse()

	// Signal for shutdown
	done := make(chan bool)

	shutdownc := make(chan io.Closer, 1)
	go findcab.HandleSignals(shutdownc)

	server := &findcab.Server{
		Port: 8080,
	}

	go func() {
		err := server.Http().ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	// Here is a list of shutdown hooks to execute when receiving the OS signal
	shutdownc <- findcab.ShutdownSequence{
		findcab.ShutdownHook(func() error {
			// Clean up database connections
			return nil
		}),
		findcab.ShutdownHook(func() error {
			done <- true
			return nil
		}),
	}

	<-done // This just blocks until a bool is sent on the channel
}

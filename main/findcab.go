package main

import (
	"flag"
	"github.com/gyokuro/findcab"
	"github.com/gyokuro/findcab/impl"
	"io"
	"log"
	"strconv"
)

// Flags from the command line
var (
	httpPort        = flag.Int("p", 8080, "http server port")
	mongoUrl        = flag.String("dbUrl", "localhost", "MongoDb url")
	mongoDbName     = flag.String("dbName", "findcab", "MongoDb database name")
	mongoCollection = flag.String("dbColl", "cabs", "MongoDb collection name")
)

func main() {

	flag.Parse()

	shutdownc := make(chan io.Closer, 1)
	go findcab.HandleSignals(shutdownc)

	// Uses the mongodb as backend datastore.
	service, err := impl.NewMongoDbCabService(*mongoUrl, *mongoDbName, *mongoCollection)
	if err != nil {
		panic(err)
	}

	httpServer := findcab.HttpServer(service)
	httpServer.Addr = ":" + strconv.Itoa(*httpPort)

	// Run the http server in a separate go routine
	// When stopping, send a true to the httpDone channel.
	// The channel done is used for getting notification on clean server shutdown.
	httpDone := make(chan bool)
	done := findcab.RunServer(httpServer, httpDone)

	log.Println("Server listening on ", httpServer.Addr)

	// Here is a list of shutdown hooks to execute when receiving the OS signal
	shutdownc <- findcab.ShutdownSequence{
		findcab.ShutdownHook(func() error {
			// Clean up database connections
			service.Close()
			return nil
		}),
		findcab.ShutdownHook(func() error {
			httpDone <- true
			return nil
		}),
	}

	<-done // This just blocks until a bool is sent on the channel
}

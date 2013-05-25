package main

import (
	"flag"
	"log"
	"os"
)

//config
const (
	LOGLEVEL = 0
)

var (
	mRedisServer *MServer
)

func main() {
	log.SetFlags(23)

	var port = flag.Int("p", 8080, "tcp listen port")
	var help = flag.Bool("h", false, "print usage info")
	flag.Parse()
	if *help {
		log.Println("multi Redis Server:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	createSharedObjects()

	mRedisServer = NewMServer("", *port)

	mRedisServer.mainLoop()
}

package main

import (
	"flag"
	"log"
	"os"
	"runtime"
)

//config
const (
	LOGLEVEL = 0 //log all that great than LOGLEVEL,other wise pass
)

var (
	mRedisServer *MServer
)

func main() {
	// sets the max number of CPUs: use all logical CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())
	//runtime.MemStats.EnableGC = false

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

package main

import (
	"log"
)

func redisLog(level int, msg ...interface{}) {
	if level < LOGLEVEL {
		return
	}
	log.Println(msg)
}

func assert(e error) {
	if e != nil {
		panic(e)
	}
}

func redisPanic(msg string) {
	panic(msg)
}

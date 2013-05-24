package main

import (
	"fmt"
)

type sharedObjectsStruct struct {
	crlf             *robj
	ok               *robj
	err              *robj
	emptybulk        *robj
	czero            *robj
	cone             *robj
	cnegone          *robj
	pong             *robj
	space            *robj
	colon            *robj
	nullbulk         *robj
	nullmultibulk    *robj
	queued           *robj
	emptymultibulk   *robj
	wrongtypeerr     *robj
	nokeyerr         *robj
	syntaxerr        *robj
	sameobjecterr    *robj
	outofrangeerr    *robj
	noscripterr      *robj
	loadingerr       *robj
	slowscripterr    *robj
	bgsaveerr        *robj
	masterdownerr    *robj
	roslaveerr       *robj
	oomerr           *robj
	plus             *robj
	messagebulk      *robj
	pmessagebulk     *robj
	subscribebulk    *robj
	unsubscribebulk  *robj
	psubscribebulk   *robj
	punsubscribebulk *robj
	del              *robj
	rpop             *robj
	lpop             *robj
	lpush            *robj
	selects          [REDIS_SHARED_SELECT_CMDS]*robj
	integers         [REDIS_SHARED_INTEGERS]*robj
	mbulkhdr         [REDIS_SHARED_BULKHDR_LEN]*robj // "<value>\r\n"
	bulkhdr          [REDIS_SHARED_BULKHDR_LEN]*robj // "$<value>\r\n"
}

var (
	shared sharedObjectsStruct
)

func createSharedObjects() {
	shared.crlf = createObject(REDIS_STRING, "\r\n")
	shared.ok = createObject(REDIS_STRING, "+OK\r\n")
	shared.err = createObject(REDIS_STRING, "-ERR\r\n")
	shared.emptybulk = createObject(REDIS_STRING, "$0\r\n\r\n")
	shared.czero = createObject(REDIS_STRING, ":0\r\n")
	shared.cone = createObject(REDIS_STRING, ":1\r\n")
	shared.cnegone = createObject(REDIS_STRING, ":-1\r\n")
	shared.nullbulk = createObject(REDIS_STRING, "$-1\r\n")
	shared.nullmultibulk = createObject(REDIS_STRING, "*-1\r\n")
	shared.emptymultibulk = createObject(REDIS_STRING, "*0\r\n")
	shared.pong = createObject(REDIS_STRING, "+PONG\r\n")
	shared.queued = createObject(REDIS_STRING, "+QUEUED\r\n")
	shared.wrongtypeerr = createObject(REDIS_STRING,
		"-ERR Operation against a key holding the wrong kind of value\r\n")
	shared.nokeyerr = createObject(REDIS_STRING,
		"-ERR no such key\r\n")
	shared.syntaxerr = createObject(REDIS_STRING,
		"-ERR syntax error\r\n")
	shared.sameobjecterr = createObject(REDIS_STRING,
		"-ERR source and destination objects are the same\r\n")
	shared.outofrangeerr = createObject(REDIS_STRING,
		"-ERR index out of range\r\n")
	shared.noscripterr = createObject(REDIS_STRING,
		"-NOSCRIPT No matching script. Please use EVAL.\r\n")
	shared.loadingerr = createObject(REDIS_STRING,
		"-LOADING Redis is loading the dataset in memory\r\n")
	shared.slowscripterr = createObject(REDIS_STRING,
		"-BUSY Redis is busy running a script. You can only call SCRIPT KILL or SHUTDOWN NOSAVE.\r\n")
	shared.masterdownerr = createObject(REDIS_STRING,
		"-MASTERDOWN Link with MASTER is down and slave-serve-stale-data is set to 'no'.\r\n")
	shared.bgsaveerr = createObject(REDIS_STRING,
		"-MISCONF Redis is configured to save RDB snapshots, but is currently not able to persist on disk. Commands that may modify the data set are disabled. Please check Redis logs for details about the error.\r\n")
	shared.roslaveerr = createObject(REDIS_STRING,
		"-READONLY You can't write against a read only slave.\r\n")
	shared.oomerr = createObject(REDIS_STRING,
		"-OOM command not allowed when used memory > 'maxmemory'.\r\n")
	shared.space = createObject(REDIS_STRING, " ")
	shared.colon = createObject(REDIS_STRING, ":")
	shared.plus = createObject(REDIS_STRING, "+")

	for j := 0; j < REDIS_SHARED_SELECT_CMDS; j++ {
		shared.selects[j] = createObject(REDIS_STRING,
			fmt.Sprintf("select %d\r\n", j))
	}
	shared.messagebulk = createStringObject("$7\r\nmessage\r\n", 13)
	shared.pmessagebulk = createStringObject("$8\r\npmessage\r\n", 14)
	shared.subscribebulk = createStringObject("$9\r\nsubscribe\r\n", 15)
	shared.unsubscribebulk = createStringObject("$11\r\nunsubscribe\r\n", 18)
	shared.psubscribebulk = createStringObject("$10\r\npsubscribe\r\n", 17)
	shared.punsubscribebulk = createStringObject("$12\r\npunsubscribe\r\n", 19)
	shared.del = createStringObject("DEL", 3)
	shared.rpop = createStringObject("RPOP", 4)
	shared.lpop = createStringObject("LPOP", 4)
	shared.lpush = createStringObject("LPUSH", 5)
	for j := 0; j < REDIS_SHARED_INTEGERS; j++ {
		shared.integers[j] = createObject(REDIS_STRING, j)
		shared.integers[j].encoding = REDIS_ENCODING_INT
	}
	for j := 0; j < REDIS_SHARED_BULKHDR_LEN; j++ {
		shared.mbulkhdr[j] = createObject(REDIS_STRING,
			fmt.Sprintf("*%d\r\n", j))
		shared.bulkhdr[j] = createObject(REDIS_STRING,
			fmt.Sprintf("$%d\r\n", j))
	}

}

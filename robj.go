package main

import (
	"container/list"
)

type robj struct {
	rtype    uint
	notused  uint
	encoding uint
	lru      uint /* lru time (relative to server.lruclock) */
	ptr      interface{}
}

func createObject(otype uint, ptr interface{}) *robj {
	o := robj{}
	o.rtype = otype
	o.encoding = REDIS_ENCODING_RAW
	o.ptr = ptr
	o.lru = 0
	return &o
}

func createStringObject(s string, length int) *robj {
	//the length is useless,for compatible with c
	return createObject(REDIS_STRING, s)
}

func createListObject() *robj {
	o := createObject(REDIS_LIST, list.New())
	o.encoding = REDIS_ENCODING_LINKEDLIST
	return o
}

func checkType(c *redisClient, o *robj, rtype int) int {
	if o.rtype != uint(rtype) {
		c.addReply(shared.wrongtypeerr)
		return REDIS_ERR
	}
	return REDIS_OK
}

package main

import (
	"strconv"
	"time"
)

type Db struct {
	dict map[interface{}]interface{}
	id   int
}

func NewDb(id int) *Db {
	return &Db{
		dict: make(map[interface{}]interface{}),
		id:   id,
	}
}

func (db *Db) set(key, value interface{}) int {
	db.dict[key] = value
	return REDIS_OK
}

func (db *Db) get(key interface{}) interface{} {
	if value, present := db.dict[key]; present {
		return value
	}
	return nil
}

func (db *Db) randomKey() interface{} {
	//use the random feature of the golang map
	for key, _ := range db.dict {
		return key
	}

	return nil
}

func noPreloadGetKeys(cmd *redisCommand, argv **robj, argc int, numkeys *int, flags int) *int {
	a := 0
	return &a
}
func zunionInterGetKeys(cmd *redisCommand, argv **robj, argc int, numkeys *int, flags int) *int {
	a := 0
	return &a
}
func renameGetKeys(cmd *redisCommand, argv **robj, argc int, numkeys *int, flags int) *int {
	a := 0
	return &a
}

/* This is the generic command implementation for EXPIRE, PEXPIRE, EXPIREAT
 * and PEXPIREAT. Because the commad second argument may be relative or absolute
 * the "basetime" argument is used to signal what the base time is (either 0
 * for *AT variants of the command, or the current time for relative expires).
 *
 * unit is either UNIT_SECONDS or UNIT_MILLISECONDS, and is only used for
 * the argv[2] parameter. The basetime is always specified in milliseconds. */

func expireGenericCommand(c *redisClient, basetime int64, unit int) {
	key := c.argv[1]

	delta, err := strconv.ParseInt(string(c.argv[2]), 10, 64)
	if err != nil {
		c.addReplyError("value is not an integer or out of range")
		return
	}

	if unit == UNIT_SECONDS {
		delta *= 1000
	}

	when := basetime + delta

	_, present := c.db.dict[string(key)]
	if !present {
		c.addReply(shared.czero)
		return
	}

	if when <= mstime() {
		delete(c.db.dict, string(key))
	} else {
		go func() {
			fakec := make(chan bool, 1)
			select {
			case <-fakec:
				//never come here
			case <-time.After(time.Millisecond * time.Duration(delta)):
				delete(c.db.dict, string(key))
			}
		}()
	}
	c.addReply(shared.cone)
}

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

const (
	/* Dup object types to RDB object types. Only reason is readability (are we
	 * dealing with RDB types or with in-memory object types?). */
	REDIS_RDB_TYPE_STRING = 0
	REDIS_RDB_TYPE_LIST   = 1
	REDIS_RDB_TYPE_SET    = 2
	REDIS_RDB_TYPE_ZSET   = 3
	REDIS_RDB_TYPE_HASH   = 4

	/* Object types for encoded objects. */
	REDIS_RDB_TYPE_HASH_ZIPMAP  = 9
	REDIS_RDB_TYPE_LIST_ZIPLIST = 10
	REDIS_RDB_TYPE_SET_INTSET   = 11
	REDIS_RDB_TYPE_ZSET_ZIPLIST = 12
	REDIS_RDB_TYPE_HASH_ZIPLIST = 13

	/* Sdecial RDB opcodes (saved/loaded with rdbSaveType/rdbLoadType). */
	REDIS_RDB_OPCODE_EXPIRETIME_MS = 252
	REDIS_RDB_OPCODE_EXPIRETIME    = 253
	REDIS_RDB_OPCODE_SELECTDB      = 254
	REDIS_RDB_OPCODE_EOF           = 255
)

func rdbIsObjectType(t int) bool {
	/* Test if a type is an object type. */
	return ((t >= 0 && t <= 4) || (t >= 9 && t <= 13))
}
func rdbWriteRaw(rdb *bufio.Writer, p []byte, l int) int {
	_, err := rdb.Write(p[:l])
	if err != nil {
		//TODO what is the err?
		return REDIS_ERR
	}
	return REDIS_OK
}

func rdbSaveStringObject(rdb *bufio.Writer, o *robj) int {
	if o.encoding == REDIS_ENCODING_INT {
	} else {

	}
	return REDIS_OK
}

func rdbSaveObject(rdb *bufio.Writer, o *robj) int {
	switch o.rtype {
	case REDIS_STRING:
		rdbSaveStringObject(rdb, o)
	}
	return REDIS_OK
}

func rdbSaveObjectType(rdb *bufio.Writer, o *robj) int {
	switch o.rtype {
	case REDIS_STRING:
		return rdbSaveType(rdb, REDIS_RDB_TYPE_STRING)
		//TODO other types
	default:
	}
	return REDIS_OK
}

func rdbSaveType(rdb *bufio.Writer, rtype uint8) int {
	return rdbWriteRaw(rdb, []byte{rtype}, 1)
}

func rdbSaveLen(rdb *bufio.Writer, length uint32) int {
	var (
		buf []byte
	)

	if length < (1 << 6) {
		/* Save a 6 bit length */
		buf[0] = byte((length & 0xFF) | (REDIS_RDB_6BITLEN << 6))
		if rdbWriteRaw(rdb, buf, 1) != REDIS_OK {
			return REDIS_ERR
		}
	} else if length < (1 << 14) {
		/* Save a 14 bit length */
		buf[0] = byte(((length >> 8) & 0xFF) | (REDIS_RDB_14BITLEN << 6))
		buf[1] = byte(length & 0xFF)

		if rdbWriteRaw(rdb, buf, 2) != REDIS_OK {
			return -1
		}
	} else {
		/* Save a 32 bit length */
		buf[0] = (REDIS_RDB_32BITLEN << 6)
		if rdbWriteRaw(rdb, buf, 1) == -1 {
			return -1
		}
		//length = htonl(length);
		//TODO length
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, length)
		if err != nil {
			redisLog(REDIS_ERR, err)
			return REDIS_ERR
		}

		if rdbWriteRaw(rdb, buf.Bytes(), 4) == -4 {
			return -1
		}
	}
	return REDIS_OK
}

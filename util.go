package main

import (
	"strconv"
	"time"
)

func bytes2ll(b []byte) (num int64, idx int) {
	for _, i := range b {
		if i < '0' || i > '9' {
			break
		}
		idx++
	}
	num, _ = strconv.ParseInt(string(b[:idx]), 10, 64)
	return
}

func ll2string(b []byte, ll int64) int {
	return copy(b, []byte(strconv.FormatInt(ll, 10)))
}

func mstime() int64 {
	//1 second = 1000 milisecond = 1000 * 1000 microsecond = 1000 * 1000 * 1000 nanosecond
	//time.Now().UnixNano() return nanoseconds
	return time.Now().UnixNano() / 1000 / 1000
}

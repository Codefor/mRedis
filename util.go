package main

import (
	"strconv"
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

#!/bin/bash

TARGET=mRedis

SOURCES='server.go client.go robj.go shared.go protocol.go util.go interface.go db.go const.go main.go log.go t_list.go t_string.go rdb.go'

echo formatting...
go fmt $SOURCES

echo building...
`go build -o $TARGET -p 4 $SOURCES`

if [ $? != 0 ]
then
    echo build failed
    exit -1
else
    echo build success
fi

if [ $# -gt 0 ]
then
    if [ $1 = 'r' ] || [ $1 = 'run' ]
    then
	shift
	echo trying to run...
	./$TARGET $@
    fi
fi

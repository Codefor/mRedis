#!/bin/bash
ALL='server.go client.go robj.go shared.go protocol.go util.go interface.go db.go const.go main.go'

echo formatting...
for i in $ALL
    do
        go fmt $i
    done

echo building...
BUILD=`go build $ALL`
if [ $? != 0 ]
then
    echo build failed
    exit -1
fi

echo trying to run...
./server

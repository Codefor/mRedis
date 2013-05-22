package main

type robj struct{
    rtype uint
    encoding uint
    lru uint        /* lru time (relative to server.lruclock) */
    refcount int
    ptr interface{}
}

package main

import (
	"bufio"
	"bytes"
    "errors"
	"net"
    "log"
)

const (
	REDIS_IOBUF_LEN     = 1024 * 16
	REDIS_REQ_INLINE    = 1
	REDIS_REQ_MULTIBULK = 2
)

type redisCommandProc func(c *redisClient)

type redisGetKeysProc func(cmd *redisCommand, argv *robj, argc int, numkeys *int, flags int) *int

type redisCommand struct {
	name   string
	proc   *redisCommandProc
	arity  int
	sflags []byte /* Flags as string represenation, one char per flag. */
	flags  int    /* The actual flags, obtained from the 'sflags' field. */
	/* Use a function to determine keys arguments in a command line. */
	getkeys_proc *redisGetKeysProc
	/* What keys should be loaded in background when calling this command? */
	firstkey     int /* The first argument that's a key (0 = no keys) */
	lastkey      int /* THe last argument that's a key */
	keystep      int /* The step between first and last key */
	microseconds int64
	calls        int64
}

type redisClient struct {
	c    *net.TCPConn
	bufr *bufio.Reader
	bufw *bufio.Writer
	db   *Db
	argv [][]byte

	qbuf []byte //query input buffer

	qbufi   int //query buffer idx
	qbufl   int //query buffer len
	reqtype int //default 0
    multibulklen int64
    bulklen int64
	lastErr error
}

func (c *redisClient) addReply(o *robj) {
	c.bufw.Write([]byte("-ERR " + "hello" +  "\r\n"))
	c.bufw.Flush()
}

func (c *redisClient) addReplyError(msg string) {
	c.bufw.Write([]byte("-ERR " + msg + "\r\n"))
	c.bufw.Flush()
}

func (c *redisClient) processCommand() {
	if string(c.argv[0]) == "quit" {
		c.addReply(shared.ok)
	}
    //log.Println("argv:",c.argv)
	c.addReply(shared.ok)
}

func (c *redisClient) processInlineBuffer() int {
	newline := bytes.IndexByte(c.qbuf[c.qbufi:], lf_byte)
	if newline == -1 {
		c.addReplyError("Protocol error: too big inline request")
		return -1
	}

	c.qbufi += newline
	for _, bf := range bytes.Split(c.qbuf[:newline],[]byte(" ")) {
		if len(bf) > 0 {
			c.argv = append(c.argv, bf)
		}
	}
	return 0
}

func (c *redisClient) processMultibulkBuffer() int {
    var idx int
    if c.multibulklen  == 0{
        //skip *
        c.qbufi += 1

        newline := bytes.IndexByte(c.qbuf[c.qbufi:], lf_byte)
        if newline == -1 {
            c.addReplyError("Protocol error: too big inline request")
            return -1
        }

        if newline > len(c.qbuf[c.qbufi:c.qbufl]) {
            c.addReplyError("error")
            return -1
        }
        //log.Println(c.qbuf[c.qbufi:c.qbufl],c.qbufi,c.qbuf[c.qbufi],c.qbuf[c.qbufi] != count_byte,count_byte)
        c.multibulklen,idx = parseInt(c.qbuf[c.qbufi:])
        //123\r\nabc
        c.qbufi += idx+2
    }
    //log.Println(c.multibulklen,string(c.qbuf[c.qbufi:c.qbufl]),c.qbuf[c.qbufi:c.qbufl])

    for c.multibulklen >0{
        //log.Println(c.multibulklen,c.bulklen)
        //$2\r\nab
        if c.bulklen == -1{
            //log.Println(c.qbuf[c.qbufi:c.qbufl],lf_byte,bytes.IndexByte(c.qbuf[c.qbufi:c.qbufl],lf_byte))
            newline := bytes.IndexByte(c.qbuf[c.qbufi:c.qbufl],lf_byte)
            if newline == -1 {
                c.lastErr = errors.New("some error")
                break
            }

            if c.qbuf[c.qbufi] != size_byte{
                c.lastErr = errors.New("Protocol error: expected '$',found:" + string(c.qbuf[c.qbufi]))
                return -1
            }

            c.bulklen,idx = parseInt(c.qbuf[c.qbufi+1:])
            c.qbufi += idx+3
            //log.Println(c.bulklen,string(c.qbuf[c.qbufi:c.qbufl]))
        }
        //log.Println(c.qbuf[c.qbufi:int64(c.qbufi)+c.bulklen])
        c.argv = append(c.argv,c.qbuf[c.qbufi:int64(c.qbufi)+c.bulklen])
        c.qbufi += int(c.bulklen) + 2

        //log.Println(c.argv)
        c.bulklen = -1
        c.multibulklen --
    }

    return REDIS_OK
}

func (c *redisClient) close() {
    c.c.Close()
}
func (c *redisClient) readQuery() {
	/**
	   * func (b *Reader) Read(p []byte) (n int, err error)
	  //Read reads data into p. It returns the number of bytes read into p. It calls Read at most once on the underlying Reader, hence n may be less than len(p). At EOF, the count will be zero and err will be io.EOF.

	  //see http://golang.org/pkg/bufio/#Reader.Read
	*/
	c.qbufl, c.lastErr = c.bufr.Read(c.qbuf)
	if c.lastErr != nil {
		return
	}

    //reset query buffer index
    c.qbufi = 0
    log.Println("query buffer:",c.qbuf[:c.qbufl],string(c.qbuf[:c.qbufl]))
	for c.qbufi < c.qbufl {
		//set query type
		if c.reqtype == 0 {
			if c.qbuf[0] == count_byte {
				c.reqtype = REDIS_REQ_MULTIBULK
			} else {
				c.reqtype = REDIS_REQ_INLINE
			}
		}

		if c.reqtype == REDIS_REQ_INLINE {
			if c.processInlineBuffer() != REDIS_OK {
                c.lastErr = errors.New("process Inline buffer error")
				break
			}
		} else if c.reqtype == REDIS_REQ_MULTIBULK {
			if c.processMultibulkBuffer() != REDIS_OK {
                c.lastErr = errors.New("process Multi buffer error")
				break
			}
		} else {
			redisPanic("Unknown request type")
		}
	}
}

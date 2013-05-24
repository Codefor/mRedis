package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	//"strconv"
	"strings"
)

const (
	REDIS_REQ_INLINE    = 1
	REDIS_REQ_MULTIBULK = 2
)

type redisGetKeysProc func(cmd *redisCommand, argv **robj, argc int, numkeys *int, flags int) *int

type redisClient struct {
	c    *net.TCPConn
	bufr *bufio.Reader
	bufw *bufio.Writer
	db   *Db
	argv [][]byte

	qbuf []byte //query input buffer

	qbufi        int //query buffer idx
	qbufl        int //query buffer len
	reqtype      int //default 0
	multibulklen int64
	bulklen      int64
	lastErr      error

	cmd     *redisCommand
	lastcmd *redisCommand

	/* Response buffer */
	bufpos int
	wbuf   []byte
}

func (c *redisClient) addReplyLongLong(ll int64) {
	if ll == 0 {
		c.addReply(shared.czero)
	} else if ll == 1 {
		c.addReply(shared.cone)
	} else {
		c.addReplyLongLongWithPrefix(ll, byte(':'))
	}
}

func (c *redisClient) addReplyLongLongWithPrefix(ll int64, prefix byte) {
	var (
		buf [128]byte
		l   int
	)

	if prefix == count_byte && ll < REDIS_SHARED_BULKHDR_LEN {
		c.addReply(shared.mbulkhdr[ll])
		return
	} else if prefix == size_byte && ll < REDIS_SHARED_BULKHDR_LEN {
		c.addReply(shared.bulkhdr[ll])
		return
	}

	buf[0] = prefix
	l = ll2string(buf[1:], ll)
	buf[l+1] = byte('\r')
	buf[l+2] = byte('\n')
	c.addReplyString(buf[:l+2])
}

func (c *redisClient) addReplyString(b []byte) {
}

func (c *redisClient) addReplyBulkLen(o *robj) {
	var l int64
	if o.encoding == REDIS_ENCODING_RAW {
		l = int64(len(o.ptr.(string)))
	} else {
		l = 10
	}
	c.addReplyLongLongWithPrefix(l, '$')
}

func (c *redisClient) addReplyBulk(o *robj) {
	c.addReplyBulkLen(o)
	c.addReply(o)
	c.addReply(shared.crlf)
}

func (c *redisClient) addReply(o *robj) {
	log.Println(o)
	if o.encoding == REDIS_ENCODING_RAW {
		c.bufw.Write([]byte(o.ptr.(string)))
	} else if o.encoding == REDIS_ENCODING_INT {
		//length := ll2string(c.wbuf, o.ptr.(int64))
		//c.bufw.Write([]byte(strconv.FormatInt(o.ptr.(int64), 10) + "\r\n"))
	}
	log.Println("buffer writer,buffered:", c.bufw.Buffered(), " avai:", c.bufw.Available())
}

func (c *redisClient) addReplyError(msg string) {
	c.bufw.Write([]byte("-ERR " + msg + "\r\n"))
}

func (c *redisClient) lookupCommand(cmd []byte) *redisCommand {
	rc, present := redisCommandTable[strings.ToLower(string(cmd))]
	if present {
		return rc
	} else {
		return nil
	}
}

func (c *redisClient) processCommand() int {
	/**
	 * return REDIS_OK if the client is still alive and valid
	 * return REDIS_ERR otherwise(client is destroyed)
	 */
	if string(c.argv[0]) == "quit" {
		c.addReply(shared.ok)
		//client destroyed
		return REDIS_OK
	}

	c.cmd = c.lookupCommand(c.argv[0])
	c.lastcmd = c.cmd

	if c.cmd == nil {
		c.addReplyError(fmt.Sprintf("unknown command '%s'", c.argv[0]))
		goto END
	} else if c.cmd.arity > 0 && c.cmd.arity != len(c.argv) {
		c.addReplyError(fmt.Sprintf("wrong number of arguments for command '%s'", c.argv[0]))
		goto END
	}

	log.Println(c.cmd, c.argv)
	c.cmd.proc(c)
	log.Println("bufw:", c.bufw.Buffered(), "avai:", c.bufw.Available())
END:
	c.bufw.Flush()
	c.reset()

	return REDIS_OK
}

func (c *redisClient) reset() int {
	c.cmd = nil
	c.argv = c.argv[:0]
	c.reqtype = 0
	c.multibulklen = 0
	c.bulklen = -1
	return REDIS_OK
}

func (c *redisClient) processInlineBuffer() int {
	newline := bytes.IndexByte(c.qbuf[c.qbufi:], lf_byte)
	if newline == -1 {
		c.addReplyError("Protocol error: too big inline request")
		return -1
	}

	c.qbufi += newline
	for _, bf := range bytes.Split(c.qbuf[:newline-1], []byte(" ")) {
		if len(bf) > 0 {
			c.argv = append(c.argv, bf)
		}
	}
	return 0
}

func (c *redisClient) processMultibulkBuffer() int {
	var idx int
	if c.multibulklen == 0 {
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
		c.multibulklen, idx = bytes2ll(c.qbuf[c.qbufi:])
		//123\r\nabc
		c.qbufi += idx + 2
	}
	//log.Println(c.multibulklen,string(c.qbuf[c.qbufi:c.qbufl]),c.qbuf[c.qbufi:c.qbufl])

	for c.multibulklen > 0 {
		//log.Println(c.multibulklen,c.bulklen)
		//$2\r\nab
		if c.bulklen == -1 {
			//log.Println(c.qbuf[c.qbufi:c.qbufl],lf_byte,bytes.IndexByte(c.qbuf[c.qbufi:c.qbufl],lf_byte))
			newline := bytes.IndexByte(c.qbuf[c.qbufi:c.qbufl], lf_byte)
			if newline == -1 {
				c.lastErr = errors.New("some error")
				break
			}

			if c.qbuf[c.qbufi] != size_byte {
				c.lastErr = errors.New("Protocol error: expected '$',found:" + string(c.qbuf[c.qbufi]))
				return -1
			}

			c.bulklen, idx = bytes2ll(c.qbuf[c.qbufi+1:])
			c.qbufi += idx + 3
			//log.Println(c.bulklen,string(c.qbuf[c.qbufi:c.qbufl]))
		}
		//log.Println(c.qbuf[c.qbufi:int64(c.qbufi)+c.bulklen])
		c.argv = append(c.argv, c.qbuf[c.qbufi:int64(c.qbufi)+c.bulklen])
		c.qbufi += int(c.bulklen) + 2

		//log.Println(c.argv)
		c.bulklen = -1
		c.multibulklen--
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
	log.Println("query buffer:", c.qbuf[:c.qbufl], string(c.qbuf[:c.qbufl]))
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
		//should there be *2$3get$1a*3$3get$1a?
		break
	}
}

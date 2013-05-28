package client

import (
	"io"
	"log"
	"net"
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

// protocol's special bytes
const (
	cr_byte    byte = byte('\r')
	lf_byte         = byte('\n')
	space_byte      = byte(' ')
	err_byte        = byte('-')
	ok_byte         = byte('+')
	count_byte      = byte('*')
	size_byte       = byte('$')
	num_byte        = byte(':')
	true_byte       = byte('1')
)

type RedisClient struct{
    conn *net.Conn
    r *bufio.Reader
}

type ctlbytes []byte

var crlf_bytes ctlbytes = ctlbytes{cr_byte, lf_byte}

// ----------------------------------------------------------------------------
// Services
// ----------------------------------------------------------------------------

// Creates the byte buffer that corresponds to the specified Command and
// provided command arguments.
//
// panics on error (with redis.Error)
func createRequestBytes(cmd string, args [][]byte) string {

	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Sprintf("createRequestBytes(%s) - failed to create request buffer", cmd))
		}
	}()
	cmd_bytes := []byte(cmd)

	buffer := bytes.NewBufferString("")
	//*4\r\n
	buffer.WriteByte(count_byte)
	buffer.Write([]byte(strconv.Itoa(len(args) + 1)))
	buffer.Write(crlf_bytes)
	//$6\r\nlrange\r\n
	buffer.WriteByte(size_byte)
	buffer.Write([]byte(strconv.Itoa(len(cmd_bytes))))
	buffer.Write(crlf_bytes)
	buffer.Write(cmd_bytes)
	buffer.Write(crlf_bytes)

	for _, s := range args {
		buffer.WriteByte(size_byte)
		buffer.Write([]byte(strconv.Itoa(len(s))))
		buffer.Write(crlf_bytes)
		buffer.Write(s)
		buffer.Write(crlf_bytes)
	}

	return buffer.String()
}

func appendAndConvert(a0 string, arr ...string) [][]byte {
	sarr := make([][]byte, 1+len(arr))
	sarr[0] = []byte(a0)
	for i, v := range arr {
		sarr[i+1] = []byte(v)
	}
	return sarr
}

func(c *RedisClient)readToCRLF() []byte {
	buf, e := c.r.ReadBytes(lf_byte)
	if e != nil {
		panic(fmt.Sprintf("readToCRLF - ReadBytes", e))
	}
	return buf[0 : len(buf)-2]
}

func(c *RedisClient)getResponse() (data interface{}) {
	buf := c.readToCRLF()
	//see:http://redis.io/topics/protocol
	switch buf[0] {
	case ok_byte:
		//status reply
		data = string(buf[1:])
	case err_byte:
		//error reply
		data = string(buf[1:])
	case num_byte:
		//interger reply
		data, _ = strconv.ParseInt(string(buf[1:]), 10, 64)
	case size_byte:
		//bulk reply $5value
		size, _ := strconv.Atoi(string(buf[1:]))
		data = c.readBulkData(size)

	case count_byte:
		//multi bulk reply *1*1$2ab can be nested
		cnt, _ := strconv.Atoi(string(buf[1:]))
		data = c.readMultiBulkData(cnt)
	}
	return
}

func(c *RedisClient)readBulkData(size int) (data string) {
	bulk_buffer := make([]byte, size+2)
	io.ReadFull(c.r, bulk_buffer)
	data = string(bulk_buffer[0:len(bulk_buffer) - 2])
	return
}

func(c *RedisClient)readMultiBulkData(num int) interface{} {
	data := make([]interface{}, num)
	for i := 0; i < num; i++ {
		line := c.readToCRLF()

		if len(line) == 0 {
			return nil
		}
		switch line[0] {
		case '+':
			data[i] = string(line[1:])
		case '-':
			data[i] = string(line[1:])
		case ':':
			n, err := strconv.ParseInt(string(line[1:]), 10, 64)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			data[i] = n
		case '$':
            //If the requested value does not exist
            //the bulk reply will use the special value -1 as data length
            //see:http://redis.io/topics/protocol
			n, err := strconv.Atoi(string(line[1:]))
			if err != nil{
				log.Fatal(err)
				return nil
			}
            if n > 0{
			    data[i] = c.readBulkData(n)
            }else{
			    data[i] = nil
            }
		case '*':
			n, err := strconv.Atoi(string(line[1:]))
			if err != nil || n < 0 {
				log.Fatal(err)
				return nil
			}
			data[i] = c.readMultiBulkData(n)
		}
	}
	return data
}

func(c *RedisClient)RANDOMKEY()string{
	req := createRequestBytes("randomkey",[][]byte{})
	fmt.Fprintf(*c.conn, req)
    return c.getResponse().(string)
}

func(c *RedisClient)SELECT(db int) string {
	req := createRequestBytes("select", [][]byte{[]byte(strconv.Itoa(db))})
	fmt.Fprintf(*c.conn, req)
    return c.getResponse().(string)
}

func(c *RedisClient)SET(key,value string) string {
	req := createRequestBytes("set", [][]byte{[]byte(key),[]byte(value)})
	fmt.Fprintf(*c.conn, req)
    return c.getResponse().(string)
}

func(c *RedisClient)GET(key string) string {
	req := createRequestBytes("get", [][]byte{[]byte(key)})
	fmt.Fprintf(*c.conn, req)
    return c.getResponse().(string)
}

func(c *RedisClient)APPEND(key string,toappend string) int64{
	req := createRequestBytes("append", [][]byte{[]byte(key),[]byte(toappend)})
	fmt.Fprintf(*c.conn, req)
    return c.getResponse().(int64)
}

func(c *RedisClient)EXISTS(key string) int64{
	req := createRequestBytes("exists", [][]byte{[]byte(key)})
	fmt.Fprintf(*c.conn, req)
    return c.getResponse().(int64)
}

func(c *RedisClient)LRANGE(key string, start int, end int)(data []interface{}){
	req := createRequestBytes("lrange", [][]byte{[]byte(key), []byte(strconv.Itoa(start)), []byte(strconv.Itoa(end))})
	fmt.Fprintf(*c.conn, req)
    data = c.getResponse().([]interface{})
    return
}

func(c *RedisClient)INFO()(data interface{}){
    req := createRequestBytes("info", [][]byte{})
	fmt.Fprintf(*c.conn, req)
    data = c.getResponse().(interface{})
    return
}


func(c *RedisClient)MGET(key string,arr ...string)(data interface{}){
	args := appendAndConvert(key, arr...)
    req := createRequestBytes("mget", args)
	fmt.Fprintf(*c.conn, req)
    data = c.getResponse().(interface{})
    return
}

func(c *RedisClient)SENTINEL(m string)(host string,port int){
	req := createRequestBytes("sentinel", [][]byte{[]byte(m)})
	fmt.Fprintf(*c.conn, req)
    master := c.getResponse().([]interface{})[0]

    if data,ok := master.([]interface{});ok{
        info := make(map[string]interface{})
        length := len(data)
        for i:=0;i<length;i+= 2{
            info[string(data[i].(string))] = data[i+1]
        }
        host        = info["ip"].(string)
        port,_      = strconv.Atoi(info["port"].(string))
    }
    return
}

func CreateRedisClientChans(c chan *RedisClient,host string,port int){
    c <- CreateRedisClient(host,port)
}

func CreateRedisClient(host string,port int)*RedisClient{
	conn, err := net.Dial("tcp", host + ":" + strconv.Itoa(port))
    if err != nil{
        log.Println(err)
    }
	r := bufio.NewReader(conn)
    return &RedisClient{&conn,r}
}

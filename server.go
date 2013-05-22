/**
* a multi thread redis like server
* @author: yugaohe@dangdang.com
* @since:2013-05-18
*/

package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	REDIS_OK = 0
)

var (
	maxRead  = 1100
	msgStop  = []byte("cmdStop")
	msgStart = []byte("cmdContinue")
)

var shared sharedObjectsStruct

type Db map[interface{}]interface{}

type MServer struct {
	host     string
	port     int
	dbCnt    int
	db       map[int]*Db
	tcpAddr  *net.TCPAddr
	listener *net.TCPListener
	timeout  time.Duration
}

func NewMServer(host string, port int) *MServer {
	s := MServer{
		host: host,
		port: port,
	}

	if s.dbCnt <= 0 {
		s.dbCnt = 16
	}

	s.db = make(map[int]*Db, s.dbCnt)
	s.tcpAddr = s.getTCPAddr()
	s.listener, _ = net.ListenTCP("tcp", s.tcpAddr)
	s.timeout = 20 * time.Second

	return &s
}

func (m *MServer) NewClient(c *net.TCPConn) *redisClient {
	client := &redisClient{
		c:    c,
		bufr: bufio.NewReader(c),
		bufw: bufio.NewWriter(c),
		qbuf: make([]byte, 1024*16),
        bulklen:-1,
	}

	m.selectDB(client, 0)
	return client
}

func (m *MServer) selectDB(c *redisClient, idx int) {
	if idx > 0 && idx < m.dbCnt {
		c.db = m.db[idx]
	}
}

func (m *MServer) getTCPAddr() (tcpAddr *net.TCPAddr) {
	lAddr := net.JoinHostPort(m.host, strconv.Itoa(m.port))
	tcpAddr, _ = net.ResolveTCPAddr("tcp", lAddr)
    return
}

func (m *MServer) mainLoop() {
	for {
		c, _ := m.listener.AcceptTCP()
		go m.process(c)
	}
}

func (m *MServer) process(c *net.TCPConn) {
	remoteAddr := c.RemoteAddr().String()
	log.Println(remoteAddr + " connected")
	client := m.NewClient(c)

	for {
		t := make(chan bool, 1)

		go func() {
			client.readQuery()
			t <- true
		}()

		select {
            case <-t:
                if client.lastErr != nil {
                    goto DISCONNECT
                }
                client.processCommand()
            case <-time.After(m.timeout):
                goto DISCONNECT
		}
	}
DISCONNECT:
    log.Println(client.lastErr)
	client.close()
	log.Println(remoteAddr + " disconnected")
}

/**
func (m *MServer) readBulkData(bufr *bufio.Reader, size int) string {
	bulk_buffer := make([]byte, size+2)
	io.ReadFull(c.r, bulk_buffer)
	data = string(bulk_buffer[0 : len(bulk_buffer)-2])
	return
}
*/

func main() {
	log.SetFlags(23)

	var port = flag.Int("p", 8080, "tcp listen port")
	var help = flag.Bool("h", false, "print usage info")
	flag.Parse()
	if *help {
		log.Println("multi Redis Server:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	mRedisServer := NewMServer("", *port)

	mRedisServer.mainLoop()
}

func assert(e error) {
	if e != nil {
		panic(e)
	}
}
func redisPanic(msg string) {
	panic(msg)
}

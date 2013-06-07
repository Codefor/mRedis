/**
* a multi thread redis like server
* @author: yugaohe@dangdang.com
* @since:2013-05-18
 */

package main

import (
	"bufio"
	"bytes"
	"fmt"
	//	"io"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	maxRead  = 1100
	msgStop  = []byte("cmdStop")
	msgStart = []byte("cmdContinue")
)

type MServer struct {
	host     string
	port     int
	dbCnt    int
	db       map[int]*Db
	tcpAddr  *net.TCPAddr
	listener *net.TCPListener
	timeout  time.Duration
	clients  []*redisClient
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

	for idx := 0; idx < s.dbCnt; idx++ {
		s.db[idx] = NewDb(idx)
	}

	s.tcpAddr = s.getTCPAddr()
	s.listener, _ = net.ListenTCP("tcp", s.tcpAddr)
	s.timeout = 6 * 3600 * time.Second

	return &s
}

func (m *MServer) NewClient(c *net.TCPConn) *redisClient {
	client := &redisClient{
		c:       c,
		bufr:    bufio.NewReader(c),
		bufw:    bufio.NewWriter(c),
		qbuf:    make([]byte, 1024*16),
		bulklen: -1,
		bufpos:  0,
		//wbuf:    make([]byte, 1024*16),
	}

	//default select db 0
	client.selectDB(0)
	return client
}

func (m *MServer) getTCPAddr() (tcpAddr *net.TCPAddr) {
	lAddr := net.JoinHostPort(m.host, strconv.Itoa(m.port))
	tcpAddr, _ = net.ResolveTCPAddr("tcp", lAddr)
	return
}

func (m *MServer) mainLoop() {
	for {
		c, _ := m.listener.AcceptTCP()

		client := m.NewClient(c)
		//record the clients,to show in `client list`
		m.clients = append(m.clients, client)

		go m.process(client)
	}
}

func (m *MServer) process(client *redisClient) {
	remoteAddr := client.c.RemoteAddr().String()
	redisLog(REDIS_NOTICE, remoteAddr+" connected")

	for {
		redisLog(REDIS_NOTICE, "start another readQuery")
		t := make(chan bool, 1)

		go func() {
			client.readQuery()
			t <- true
		}()

		select {
		case <-t:
			if client.lastErr != nil {
				redisLog(REDIS_NOTICE, "readQuery encountered error", client.lastErr)
				goto DISCONNECT
			}
			redisLog(REDIS_NOTICE, "start process query")
			client.processCommand()
			redisLog(REDIS_NOTICE, "end process query")
		case <-time.After(m.timeout):
			redisLog(REDIS_NOTICE, "readQuery timeout")
			goto DISCONNECT
		}
		redisLog(REDIS_NOTICE, "end another readQuery")
	}
DISCONNECT:
	redisLog(REDIS_NOTICE, client.lastErr)
	client.close()
	redisLog(REDIS_NOTICE, remoteAddr+" disconnected")
}

func (m *MServer) rdbSave(rdb_filename string) int {
	var (
		magic []byte
		err   error
	)

	tmpFile := fmt.Sprintf("temp-%d.rdb", os.Getpid())

	f_out_rdb, err := os.Open(tmpFile)
	if err != nil {
		redisLog(REDIS_WARNING, err)
		return REDIS_ERR
	}
	defer f_out_rdb.Close()
	fp := bufio.NewWriter(f_out_rdb)

	copy(magic, fmt.Sprintf("REDIS%04d", REDIS_RDB_VERSION))
	if rdbWriteRaw(fp, magic, len(magic)) != REDIS_OK {
		goto werr
	}

	for idx := 0; idx < m.dbCnt; idx++ {
		for key, value := range m.db[idx].dict {
			//save SELECT DB opcode
			if rdbSaveType(fp, REDIS_RDB_OPCODE_SELECTDB) != REDIS_OK {
				goto werr
			}
			if rdbSaveLen(fp, uint32(idx)) != REDIS_OK {
				goto werr
			}
			redisLog(REDIS_DEBUG, key, value)
		}
	}
werr:
	return REDIS_ERR
}

func getAllClientsInfoString() []byte {
	s := new(bytes.Buffer)
	for _, c := range mRedisServer.clients {
		s.Write([]byte(getClientInfoString(c)))
		s.Write([]byte("\n"))
	}

	return s.Bytes()
}

func getClientInfoString(c *redisClient) string {
	return "Addr:" + c.c.RemoteAddr().String()
}

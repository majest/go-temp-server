package server

import (
	"bytes"
	"fmt"
	log "github.com/cihub/seelog"
	m "github.com/majest/sambo-go-tcp-server/mysql"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	RECV_BUF_LEN = 32
)

type callback func(string)

type Server struct {
	listener net.Listener
	mysql    *m.Db
	Call     callback
	test     bool
	db       *m.Db
	buffer   map[string]*bytes.Buffer
	lock     sync.Mutex
}

func New(port string, test bool) *Server {

	if test {
		log.Debug("Starting the server in test mode")
	} else {
		log.Debug("Starting the server")
	}
	server := &Server{}
	server.test = test
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Errorf("error listening:", err.Error())
		os.Exit(1)
	}
	server.listener = listener
	server.buffer = make(map[string]*bytes.Buffer)
	return server
}

func (s *Server) SetDb(db *m.Db) {
	s.db = db
}

func (s *Server) Process() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Errorf("Error accept:", err.Error())
			return
		}

		go s.Receive(conn)
	}
}

func (s *Server) getBuffer(location string) *bytes.Buffer {

	// create buffer if it's not there
	if _, ok := s.buffer[location]; !ok {
		s.buffer[location] = &bytes.Buffer{}
	}

	return s.buffer[location]
}

func (s *Server) CheckAndSave(location string) {

	// lock the method
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.test && s.db != nil {

		log.Debugf("Checking data for %s", location)
		buffer := s.getBuffer(location)
		data := buffer.String()

		if strings.Contains(data, "$") {

			log.Debugf("Processing and Saving data: %s", data)
			// split by $
			parts := strings.Split(data, "$")

			//log.Debugf("========= %v", []byte(parts[0]))
			// add ip and insert data
			s.db.InsertData(fmt.Sprintf("%s;%s", parts[0], location))

			// clear buffer
			buffer.Truncate(0)

			// add remaining part of the packet back tu buffer
			buffer.WriteString(parts[1])
		}
	}
}

func clear(b []byte) []byte {
	r := []byte{}

	for _, v := range b {
		if v != 0 {
			r = append(r, v)
		}
	}
	return r
}

func (s *Server) Receive(conn net.Conn) {

	for {
		buf := make([]byte, RECV_BUF_LEN)
		n, err := conn.Read(buf)
		if err != nil {
			log.Debugf("Closing connection")
			conn.Close()
			return
		}

		log.Debugf("received %v bytes. Data: %s", n, string(buf))

		// get the remote address
		location := conn.RemoteAddr().String()
		s.getBuffer(location).WriteString(string(clear(buf)))
		s.CheckAndSave(location)
	}
}

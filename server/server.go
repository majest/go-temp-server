package server

import (
	"bytes"
	"fmt"
	log "github.com/cihub/seelog"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	RECV_BUF_LEN = 32
)

type Storage interface {
	Save(string, string)
}

type callback func(string)

type Server struct {
	listener net.Listener
	buffer   map[string]*bytes.Buffer
	lock     sync.Mutex
	storage  Storage
	port     string
}

func New(port string, test bool) *Server {

	if test {
		log.Debug("Starting the server in test mode")
	} else {
		log.Debug("Starting the server")
	}

	server := &Server{}
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Errorf("error listening:", err.Error())
		os.Exit(1)
	}

	server.listener = listener
	server.buffer = make(map[string]*bytes.Buffer)
	server.port = port
	return server
}

func (s *Server) Process(storage Storage) {
	s.storage = storage

	log.Infof("Started TCP listener on port %v", s.port)
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

func (s *Server) CheckAndSave(location string, storage Storage) {

	// lock the method
	s.lock.Lock()
	defer s.lock.Unlock()

	log.Debugf("Checking data for %s", location)
	buffer := s.getBuffer(location)
	data := buffer.String()

	if strings.Contains(data, "\n") {

		//	log.Debugf("Processing and Saving data: %s", data)
		// split by $
		parts := strings.Split(data, "\n")

		//log.Debugf("========= %v", []byte(parts[0]))
		// add ip and insert data
		log.Infof("got message with CR")
		storage.Save(parts[0], location)

		// clear buffer
		buffer.Truncate(0)

		// add remaining part of the packet back to buffer
		buffer.WriteString(parts[1])
	} else {
		log.Infof("got message without CR")
		buffer.WriteString(data)
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
		s.CheckAndSave(location, s.storage)
	}
}

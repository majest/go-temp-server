package server

import (
	"fmt"
	log "github.com/cihub/seelog"
	m "github.com/majest/sambo-go-tcp-server/mysql"
	"net"
	"os"
)

const (
	RECV_BUF_LEN = 2048
)

type callback func(string)

type Server struct {
	listener net.Listener
	mysql    *m.Db
	Call     callback
	test     bool
	db       *m.Db
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

func (s *Server) Receive(conn net.Conn) {

	for {
		buf := make([]byte, RECV_BUF_LEN)
		n, err := conn.Read(buf)
		if err != nil {
			log.Debugf("Closing connection")
			conn.Close()
			return
		}

		data := string(buf)

		log.Debugf("received %s bytes of data = %s", n, data)

		if !s.test && s.db != nil {
			s.db.InsertData(data)
		}
	}
}

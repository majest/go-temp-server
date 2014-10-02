package main

import (
	"flag"
	s "github.com/majest/go-temp-server/server"
	storage "github.com/majest/go-temp-server/storage/file"
)

var port string
var mysqlHost string
var mysqlPort string
var mysqlUser string
var mysqlPassword string
var mysqlDbName string
var mysqlTable string
var testMode bool

func init() {
	flag.StringVar(&port, "port", "9002", "Server port")
}

func main() {
	flag.Parse()

	// init storage
	st := storage.New()

	// server
	server := s.New(port, false)

	// init processing
	go server.Process(st)

	// start storage
	st.Run()
}

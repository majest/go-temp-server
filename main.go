package main

import (
	"flag"
	m "github.com/majest/sambo-go-tcp-server/mysql"
	s "github.com/majest/sambo-go-tcp-server/server"
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
	flag.StringVar(&port, "port", "5566", "Server port")
	flag.StringVar(&mysqlHost, "mysql-host", "127.0.0.1", "MySql host/ip address")
	flag.StringVar(&mysqlPort, "mysql-port", "3306", "MySql port")
	flag.StringVar(&mysqlUser, "mysql-user", "root", "MySql user")
	flag.StringVar(&mysqlPassword, "mysql-password", "", "MySql user password")
	flag.StringVar(&mysqlDbName, "mysql-dbname", "wagi", "Database name")
	flag.StringVar(&mysqlTable, "mysql-table", "wagi", "Table name")
	flag.BoolVar(&testMode, "testMode", false, "Check the server without any database interaction")

}

func main() {
	flag.Parse()

	server := s.New(port, testMode)

	if !testMode {
		db := m.New(mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDbName, mysqlTable)
		server.SetDb(db)
	}

	server.Process()
}

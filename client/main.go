package main

import (
	"net"
	"os"
	"time"
)

func main() {
	strEcho := "Halo\n"

	servAddr := "localhost:9002"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	i := 0
	for {

		_, err = conn.Write([]byte(strEcho))
		if err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		println("write to server = ", strEcho)

		if i%5 == 0 {
			_, err = conn.Write([]byte("\n"))
			if err != nil {
				println("Write to server failed:", err.Error())
				os.Exit(1)
			}
		}

		i++
		time.Sleep(3 * time.Second)
	}

	conn.Close()
}

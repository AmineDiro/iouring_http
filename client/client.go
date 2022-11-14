package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const (
	HOST = "localhost"
	PORT = "8000"
	TYPE = "tcp"
)

var (
	ip          = flag.String("ip", "localhost", "server IP")
	connections = flag.Int("conn", 1, "number of websocket connections")
)

func main() {

	flag.Usage = func() {
		io.WriteString(os.Stderr, `Websockets client generator
	Example usage: ./client -ip=172.17.0.1 -conn=10
	`)
		flag.PrintDefaults()
	}
	flag.Parse()

	tcpServer, err := net.ResolveTCPAddr(TYPE, *ip+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	log.Printf("Connecting to %s", tcpServer)
	var conns []net.Conn
	for i := 0; i < *connections; i++ {
		c, err := net.DialTCP(TYPE, nil, tcpServer)
		if err != nil {
			println("Dial failed:", err.Error())
			os.Exit(1)
		}
		conns = append(conns, c)
		defer func() {
			c.Close()
		}()
	}

	log.Printf("Finished initializing %d connections", len(conns))

	tts := time.Second
	if *connections > 100 {
		tts = time.Millisecond * 5
	}
	for {
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			_, err = conn.Write([]byte("This is a message"))
			if err != nil {
				println("Write data failed:", err.Error())
				os.Exit(1)
			}

			time.Sleep(1 * time.Second)
		}
	}

}

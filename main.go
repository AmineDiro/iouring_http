package main

import (
	"fmt"
	"net"
	_ "net/http/pprof"
	"time"

	"os"
	"os/signal"
	"syscall"

	ringListener "github.com/aminediro/iouring_server/server"
)

func handler(conn net.Conn) {
	b := make([]byte, 512)
	for {
		_, _ = conn.Read(b)
		// conn.Close()
		time.Sleep(time.Millisecond * 10)
	}
}

func serve(l net.Listener) {
	for {
		conn, _ := l.Accept()
		go handler(conn)
	}
}

func main() {
	// runtime.GOMAXPROCS(2)

	ringListener.IncreaseResources()
	sigChannel := make(chan os.Signal, 1)

	l, err := ringListener.MKRingListener(":8000")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening ring %v\n", l)
	go l.Listen()
	go serve(l)

	// Enable pprof hooks
	// go func() {
	// 	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
	// 		log.Fatalf("pprof failed: %v", err)
	// 	}
	// }()

	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)
	<-sigChannel

	l.Close()

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	// The "/" pattern matches everything, so we need to check
	// 	// that we're at the root here.
	// 	if req.URL.Path != "/" {
	// 		http.NotFound(w, req)
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "hello io_uring!\n")
	// })

	// s := http.Server{Handler: mux}
	// if err := s.Serve(l); err != nil {
	// 	log.Fatal(err)
	// }
}

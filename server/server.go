package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"unsafe"
)

//#cgo LDFLAGS: -luring
//#include "server.h"
//#include <liburing.h>
import "C"

type RINGOP int

const (
	ACCEPT RINGOP = 1 + iota
	READ
	WRITE
)

const (
	MAXCONN     int = 10000
	SOReuseport int = 0x0F
)

var (
	connChan chan RingConn
)

type RingListener struct {
	socketFD   int
	ringCtx    context.Context
	cancelFunc context.CancelFunc
}

// Init the ring and the communication channels
func Init() error {
	ret := int(C.ring_init())
	if ret != 0 {
		return fmt.Errorf("error in ring_init")
	}
	connChan = make(chan RingConn, MAXCONN)
	return nil
}
func MKRingListener(addr string) (*RingListener, error) {
	// Create The TCP Socket + Listen syscall
	socketFD, err := bindTCPSocket(addr)
	if err != nil {
		return nil, err
	}
	// Creating Ring, setting up communication channels
	if err = Init(); err != nil {
		return nil, err
	}

	ringCtx, cancel := context.WithCancel(context.Background())
	return &RingListener{
		socketFD:   socketFD,
		ringCtx:    ringCtx,
		cancelFunc: cancel,
	}, nil
}

func (rl *RingListener) Listen() {
	// Start the ring Loop in C
	go C.ring_loop(C.int(rl.socketFD))
}

// Called by the server to get the new accepted conns
func (rl *RingListener) Accept() (RingConn, error) {
	conn := <-connChan
	fmt.Printf("Accepted conn : %#v\n", conn)
	return conn, nil
}

// Closes the listener
func (rl *RingListener) Close() {
	rl.cancelFunc()
	log.Println("Closing the ring")
	C.ring_close()
}

// Gets the formated address
func (rl *RingListener) Addr() net.Addr { return nil }

//export accept_callback
func accept_callback(connFD C.int) {
	connChan <- RingConn{client_socket: int(connFD)}
}

//export read_callback
func read_callback(iovec *C.char, length C.int) {
	// TODO : Use some kind of preallocated buffer
	// from memorypool
	// readLength := int(length)
	// buff := make([]byte, readLength)
	// copy(buff, (*(*[129]byte)(unsafe.Pointer(iovec)))[:readLength:readLength])

	buff := C.GoBytes(unsafe.Pointer(iovec), length)
	fmt.Printf("Received :%s\n", buff)

}

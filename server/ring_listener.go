package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"syscall"
)

//#cgo LDFLAGS: -luring
//#include "ring.h"
//#include <liburing.h>
import "C"

type RINGOP int

const (
	ACCEPT RINGOP = 1 + iota
	READ
	WRITE
)

const (
	MAXCONN     int = 10_000
	CHANPOOLMAX int = 10
	READSIZE    int = 512
	SOReuseport int = 0x0F
)

type ConnMap struct {
	mu      *sync.RWMutex
	chanMap map[int]chan []byte
}

var (
	acceptChan chan RingConn
	chanPool   chan chan []byte
	connMap    *ConnMap
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
		return fmt.Errorf("%v", syscall.Errno(-ret))
	}
	acceptChan = make(chan RingConn, MAXCONN)

	connMap = &ConnMap{
		mu:      &sync.RWMutex{},
		chanMap: make(map[int]chan []byte, MAXCONN),
	}
	//TODO : use this later
	chanPool = make(chan chan []byte, CHANPOOLMAX)

	for i := 0; i < CHANPOOLMAX; i++ {
		chanPool <- make(chan []byte)
	}
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

// Gets the formated address
func (rl *RingListener) Addr() net.Addr { return nil }

// Start the ring Loop in C
func (rl *RingListener) Listen() {
	C.ring_loop(C.int(rl.socketFD))
}

// Called by the server to get the new accepted conns
func (rl *RingListener) Accept() (net.Conn, error) {
	select {
	case <-rl.ringCtx.Done():
		return RingConn{}, rl.ringCtx.Err()

	case conn := <-acceptChan:
		log.Printf("Accepted conn : %#v\n", conn)
		// Creating the read channel for the conn in the conn Map
		connMap.mu.Lock()
		defer connMap.mu.Unlock()

		rcChan := make(chan []byte, 10)
		if _, ok := connMap.chanMap[conn.clientSocketFD]; !ok {
			connMap.chanMap[conn.clientSocketFD] = rcChan
		}

		return conn, nil
	}
}

// Closes the listener
func (rl *RingListener) Close() error {
	rl.cancelFunc()
	log.Println("Closing the ring")
	C.ring_close()
	return nil
}

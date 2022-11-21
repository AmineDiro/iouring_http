package server

import (
	"context"
	"log"
	"sync"
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
	MAXCONN     int = 100000
	CHANPOOLMAX int = 10
	SOReuseport int = 0x0F
)

type ConnMap struct {
	mu      *sync.RWMutex
	chanMap map[int]chan []byte
}

var (
	acceptChan chan RingConn
	readChan   chan []byte
	chanPool   chan chan []byte
	connMap    *ConnMap
)

type RingListener struct {
	socketFD   int
	ringCtx    context.Context
	cancelFunc context.CancelFunc
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

// Start the ring Loop in C
func (rl *RingListener) Listen() {
	go C.ring_loop(C.int(rl.socketFD))
}

// Called by the server to get the new accepted conns
func (rl *RingListener) Accept() (RingConn, error) {
	select {
	case <-rl.ringCtx.Done():
		return RingConn{}, rl.ringCtx.Err()

	case conn := <-acceptChan:
		log.Printf("Accepted conn : %#v\n", conn)
		return conn, nil
	}
}

// Closes the listener
func (rl *RingListener) Close() {
	rl.cancelFunc()
	log.Println("Closing the ring")
	C.ring_close()
}

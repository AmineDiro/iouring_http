package server

import (
	"fmt"
	"net"
	"sync"
	"unsafe"
)

//#cgo LDFLAGS: -luring
//#include "ring.h"
//#include <liburing.h>
import "C"

// Init the ring and the communication channels
func Init() error {
	ret := int(C.ring_init())
	if ret != 0 {
		return fmt.Errorf("can't  create ring")
	}
	acceptChan = make(chan RingConn, MAXCONN)

	connMap = &ConnMap{
		mu:      &sync.RWMutex{},
		chanMap: make(map[int]chan []byte),
	}
	chanPool = make(chan chan []byte, CHANPOOLMAX)

	for i := 0; i < CHANPOOLMAX; i++ {
		chanPool <- make(chan []byte, 100)
	}
	return nil
}

// Gets the formated address
func (rl *RingListener) Addr() net.Addr { return nil }

//export accept_callback
func accept_callback(connFD C.int) {
	acceptChan <- RingConn{clientSocketFD: int(connFD)}
}

//export read_callback
func read_callback(clientFD C.int, iovec *C.char, length C.int) {
	// TODO : Use some kind of preallocated buffer
	// The C.GoBytes create a bytes buffer and copies the iovec content
	// from memorypool
	// readLength := int(length)
	// buff <- MemPool
	// copy(buff, (*(*[129]byte)(unsafe.Pointer(iovec)))[:readLength:readLength])

	connMap.mu.RLock()
	defer connMap.mu.RUnlock()

	buff := C.GoBytes(unsafe.Pointer(iovec), length)
	connMap.chanMap[int(clientFD)] <- buff
}

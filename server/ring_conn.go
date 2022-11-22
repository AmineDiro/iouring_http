package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

//#cgo LDFLAGS: -luring
//#include "ring.h"
//#include <liburing.h>
import "C"

type RingConn struct {
	clientSocketFD int
}

// Read reads data from the connection.
// The conn submits to the ring and then waits on
func (rc RingConn) Read(b []byte) (n int, err error) {
	connMap.mu.RLock()

	//Submit read request
	rcChan, ok := connMap.chanMap[rc.clientSocketFD]
	connMap.mu.RUnlock()

	if !ok {
		return -1, fmt.Errorf("can't find channel for conn in map")
	}
	b = <-rcChan
	log.Printf("FD: %d , msg : %s", rc.clientSocketFD, string(b))

	return len(b), nil
}

// Write writes data to the connection.
func (rc RingConn) Write(b []byte) (n int, err error) { return -1, nil }

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (rc RingConn) Close() error { return nil }

// LocalAddr returns the local network address, if known.
func (rc RingConn) LocalAddr() net.Addr { return nil }

// RemoteAddr returns the remote network address, if known.
func (rc RingConn) RemoteAddr() net.Addr { return nil }

// THese will be skipped  for now
func (rc RingConn) SetDeadline(t time.Time) error      { return nil }
func (rc RingConn) SetReadDeadline(t time.Time) error  { return nil }
func (rc RingConn) SetWriteDeadline(t time.Time) error { return nil }

package server

import (
	"net"
	"time"
)

type RingConn struct {
	client_socket int
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (rc *RingConn) Read(b []byte) (n int, err error) { return -1, nil }

// Write writes data to the connection.
func (rc *RingConn) Write(b []byte) (n int, err error) { return -1, nil }

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (rc *RingConn) Close() error { return nil }

// LocalAddr returns the local network address, if known.
func (rc *RingConn) LocalAddr() net.Addr { return nil }

// RemoteAddr returns the remote network address, if known.
func (rc *RingConn) RemoteAddr() net.Addr { return nil }

// THese will be skipped  for now
func (rc *RingConn) SetDeadline(t time.Time) error      { return nil }
func (rc *RingConn) SetReadDeadline(t time.Time) error  { return nil }
func (rc *RingConn) SetWriteDeadline(t time.Time) error { return nil }

package server

import (
	"fmt"
	"log"
	"net"
	"syscall"
)

const (
	MAXCONN     int = 10000
	SOReuseport int = 0x0F
)

type RingListener struct {
	socketFD int
	ring     *IOURing
	conns    chan net.Conn
}

// Increase resources limitations

func IncreaseResources() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

}

func MKRingListener(addr string) (*RingListener, error) {
	// Create The TCP Socket + Listen syscall
	socketFD, err := bindTCPSocket(addr)
	if err != nil {
		return nil, err
	}
	// Creating Ring
	ring, err := MkRing(socketFD)
	if err != nil {
		return nil, err
	}

	return &RingListener{
		socketFD: socketFD,
		ring:     ring,
		conns:    make(chan net.Conn, MAXCONN),
	}, nil
}

// Called by the server to get the new accepted conns
func (rl *RingListener) Accept() (net.Conn, error) { return <-rl.conns, nil }

// Closes the listener
func (rl *RingListener) Close() error {
	log.Println("Closing the ring")
	rl.ring.Close()
	return nil
}

// Gets the formated address
func (rl *RingListener) Addr() net.Addr { return nil }

// Start Listening by calling into the C liburing
// Either we continuously submit
func (rl *RingListener) Listen() {
	// Run the accept and printing of conns
	rl.ring.Run()
}

// Starts a socket and binds it to the address
// returns the file descriptor of the binded socket
func bindTCPSocket(address string) (int, error) {
	// Create a socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return -1, fmt.Errorf("could not open socket")
	}
	netAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return -1, fmt.Errorf("could not open socket")
	}

	var ipAddr [4]byte
	copy(ipAddr[:], netAddr.IP)
	sockAddr := &syscall.SockaddrInet4{
		Port: netAddr.Port,
		Addr: ipAddr,
	}
	// Set socket options

	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		syscall.Close(fd)
		return -1, err
	}
	// Bind socket to the adress
	if err := syscall.Bind(fd, sockAddr); err != nil {
		syscall.Close(fd)
		return -1, err
	}

	// Listen of the socket
	if err := syscall.Listen(fd, MAXCONN); err != nil {
		syscall.Close(fd)
		return -1, err
	}
	return fd, nil

}

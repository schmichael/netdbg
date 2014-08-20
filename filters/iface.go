package filters

import (
	"net"
	"strconv"
)

type Role int

const (
	Client Role = 0
	Server Role = 1
)

func (r Role) String() string {
	switch r {
	case Client:
		return "client"
	case Server:
		return "server"
	default:
		panic("unknown role: " + strconv.Itoa(int(r)))
	}
}

// FilterFactory is the New function Filters must implement to be used.
type FilterFactory func(r Role, in <-chan []byte, out chan<- []byte) (f Filter)

// Filter is the core interface for implementing network traffic filters.
type Filter interface {
	// Accept is called when a new incoming connection is created. Returning
	// false will cause the connection to be closed.
	//
	// Only server filters have Accept called.
	Accept(c net.Conn) bool

	// Close is called when the connection ends and is passed the error. If the
	// filter returns true, netdbg continues operation. If false, it exits.
	Close(err error) (ok bool)
}

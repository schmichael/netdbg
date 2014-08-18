package netdbg

import (
	"errors"
	"net"
)

type Filter interface {
	// Accept is called when a new incoming connection is created. Returning
	// false will cause the connection to be closed. Returning true will cause a
	// new outgoing connection to the target.
	Accept(c net.Conn) bool

	// Write is called when data is written by an incoming connection.
	Write(p []byte) error

	// Read is called when data is received from the target host.
	Read(p []byte) error
}

type payload struct {
	p   []byte
	err error
}

func Proxy(listener net.Listener, target string, filter Filter) error {
	for {
		incoming, err := listener.Accept()
		if err != nil {
			return err
		}
		if !filter.Accept(incoming) {
			// filter says to throw out this connection
			incoming.Close()
			continue
		}

		outgoing, err := net.Dial("tcp", target)
		if err != nil {
			incoming.Close()
			return err
		}

		in := make(chan payload)
		out := make(chan payload)
		go read(incoming, in)
		go read(outgoing, out)

		for err == nil {
			select {
			case inp := <-in:
				if inp.err != nil {
					err = inp.err
					continue
				}
				if err = filter.Write(inp.p); err != nil {
					continue
				}
				_, err = outgoing.Write(inp.p)
			case outp := <-out:
				if outp.err != nil {
					err = outp.err
					continue
				}
				if err = filter.Read(outp.p); err != nil {
					continue
				}
				_, err = incoming.Write(outp.p)
			}
		}
		incoming.Close()
		outgoing.Close()
		return err
	}
}

func read(conn net.Conn, comm chan payload) {
	for {
		buf := make([]byte, 10*1024) // this buffer size chosen using alchemy
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			comm <- payload{err: err}
			return
		}
		if n == 0 {
			conn.Close()
			comm <- payload{err: errors.New("eof i guess")}
			return
		}
		comm <- payload{p: buf[:n]}
	}
}

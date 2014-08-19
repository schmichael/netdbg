package netdbg

import (
	"errors"
	"net"

	"github.com/schmichael/netdbg/filters"
)

type payload struct {
	p   []byte
	err error
}

// Proxy routes data from clients that connect to the listener and the target.
// All traffic is passed through filters which are created per-connection.
// Connections are handled serially as are calls to the filter chain.
func Proxy(listener net.Listener, target string, filterFuncs []filters.FilterFactory) error {
	for {
		writerInput := make(chan []byte)
		readerInput := make(chan []byte)

		// Create fresh filters for each connection
		win := writerInput
		rin := readerInput
		var filter filters.Filter
		filterChain := []filters.Filter{}
		for _, filterFunc := range filterFuncs {
			var filter filters.Filter
			filter, win, rin = filterFunc(win, rin)
			filterChain = append(filterChain, filter)
		}

		// Last two returned chans are the outputs
		writerOutput := win
		readerOutput := rin

		incoming, err := listener.Accept()
		if err != nil {
			return err
		}
		for _, filter := range filterChain {
			if !filter.Accept(incoming) {
				// filter says to throw out this connection
				incoming.Close()
				continue
			}
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
				writerInput <- inp.p
			case writerOut := <-writerOutput:
				_, err = outgoing.Write(writerOut)
			case outp := <-out:
				if outp.err != nil {
					err = outp.err
					continue
				}
				readerInput <- outp.p
			case readerOut := <-readerOutput:
				_, err = incoming.Write(readerOut)
			}
		}
		close(writerInput)
		close(readerInput)
		incoming.Close()
		outgoing.Close()
		for _, filter := range filterChain {
			if !filter.Close(err) {
				// Don't return the connection error - only return unexpected internal
				// errors.
				return nil
			}
		}
	}
}

func read(conn net.Conn, comm chan payload) {
	for {
		buf := make([]byte, 10*1024) // this buffer size chosen using alchemy
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			select {
			case comm <- payload{err: err}:
			default:
				// receiver may have already died so don't block
			}
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

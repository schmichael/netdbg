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
func Proxy(listener net.Listener, target string, clientFF, serverFF []filters.FilterFactory) error {
	for {
		// Setup client chain
		clientFilters := []filters.Filter{}
		clientIn := make(chan []byte)
		clientOut := clientIn
		cin := clientIn
		for _, filterFunc := range clientFF {
			clientOut = make(chan []byte)
			filter := filterFunc(filters.Client, cin, clientOut)
			clientFilters = append(clientFilters, filter)
			cin = clientOut
		}

		// Setup server chain
		serverFilters := []filters.Filter{}
		serverIn := make(chan []byte)
		serverOut := serverIn
		sin := serverIn
		for _, filterFunc := range serverFF {
			serverOut = make(chan []byte)
			filter := filterFunc(filters.Server, sin, serverOut)
			serverFilters = append(serverFilters, filter)
			sin = serverOut
		}

		// Wait for a new connection
		client, err := listener.Accept()
		if err != nil {
			return err
		}
		for _, filter := range serverFilters {
			if !filter.Accept(client) {
				// filter says to throw out this connection
				client.Close()
				continue
			}
		}

		server, err := net.Dial("tcp", target)
		if err != nil {
			client.Close()
			return err
		}

		clientAct := make(chan payload)
		serverAct := make(chan payload)
		go read(client, clientAct)
		go read(server, serverAct)

		for err == nil {
			select {
			case act := <-clientAct:
				if act.err != nil {
					err = act.err
					continue
				}
				clientIn <- act.p
			case p := <-clientOut:
				_, err = server.Write(p)
			case act := <-serverAct:
				if act.err != nil {
					err = act.err
					continue
				}
				serverIn <- act.p
			case p := <-serverOut:
				_, err = client.Write(p)
			}
		}

		// There was an error
		close(clientIn)
		close(serverIn)
		client.Close()
		server.Close()
		for _, filterChain := range [][]filters.Filter{clientFilters, serverFilters} {
			for _, filter := range filterChain {
				if !filter.Close(err) {
					// Don't return the connection error - only return unexpected internal
					// errors.
					return nil
				}
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

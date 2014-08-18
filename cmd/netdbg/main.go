package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/schmichael/netdbg"
	"github.com/schmichael/netdbg/filters"
)

func usage() {
	fmt.Fprintf(os.Stderr, "%s usage: %s [listen host:port] [target host:port]\n", os.Args[0], os.Args[0])
	os.Exit(1)
}

func main() {
	flag.Parse()
	numArgs := len(flag.Args())

	if numArgs != 2 {
		usage()
	}
	laddrFlag := flag.Arg(0)
	target := flag.Arg(1)

	// Create listener
	laddr, err := net.ResolveTCPAddr("tcp", laddrFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid listen address %s: %v\n\n", laddrFlag, err)
		usage()
	}

	listen, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listening on %s: %v\n", laddrFlag, err)
		os.Exit(2)
	}

	// Validate target address
	_, err = net.ResolveTCPAddr("tcp", target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid target address %s: %v\n\n", target, err)
		usage()
	}

	// Create filter
	filter := filters.HumanLogger{}

	fmt.Fprintf(os.Stderr, "started %s → %v → %s\n", laddrFlag, filter, target)
	fmt.Fprintf(os.Stderr, "exited with: %v\n", netdbg.Proxy(listen, target, &filter))
}

type nopWriteCloser struct {
	io.Writer
}

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

func (nopWriteCloser) Close() error { return nil }

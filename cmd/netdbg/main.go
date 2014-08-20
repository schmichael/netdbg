package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/schmichael/netdbg"
	"github.com/schmichael/netdbg/filters"
)

const usageStr = `%s usage: %s [filters] [listen host:port] [target host:port]
filters = client-filter-1:client-filter-2,server-filter-1:server-filter-2
`

func usage() {
	fmt.Fprintf(os.Stderr, usageStr, os.Args[0], os.Args[0])
	os.Exit(1)
}

func main() {
	flag.Parse()
	numArgs := len(flag.Args())

	if numArgs != 3 {
		usage()
	}
	filterFlag := flag.Arg(0)
	laddrFlag := flag.Arg(1)
	target := flag.Arg(2)

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
	filterParts := strings.Split(filterFlag, ",")
	if len(filterParts) == 0 {
		fmt.Fprintf(os.Stderr, "missing filters\n")
		usage()
	}
	clientFilters := parseFilters(filterParts[0])
	var serverFilters []filters.FilterFactory
	if len(filterParts) > 1 {
		serverFilters = parseFilters(filterParts[1])
	}

	fmt.Fprintf(os.Stderr, "started %s ⇄ %v ⇄ %s\n", laddrFlag, filterFlag, target)
	err = netdbg.Proxy(listen, target, clientFilters, serverFilters)
	fmt.Fprintf(os.Stderr, "exited with: %v\n", err)
}

func parseFilters(filterList string) []filters.FilterFactory {
	filterChain := []filters.FilterFactory{}
	for _, name := range strings.Split(filterList, ":") {
		ffact := filters.GetFilter(name)
		if ffact == nil {
			fmt.Fprintf(os.Stderr, "unknown filter: %s\n", name)
			usage()
		}
		filterChain = append(filterChain, ffact)
	}
	return filterChain
}

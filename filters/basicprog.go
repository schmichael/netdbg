package filters

import (
	"fmt"
	"io"
	"net"
)

type BasicProgress struct{}

func (*BasicProgress) Accept(net.Conn) bool {
	fmt.Print("⇄")
	return true
}

func (*BasicProgress) Write([]byte) error {
	fmt.Print("→")
	return nil
}

func (*BasicProgress) Read([]byte) error {
	fmt.Print("←")
	return nil
}

func (*BasicProgress) Close(err error) bool {
	if err == io.EOF {
		fmt.Print("↹\n")
	} else {
		fmt.Print("↯\n")
	}
	return true
}

func (*BasicProgress) String() string {
	return "basic-progress"
}

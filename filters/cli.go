package filters

import (
	"fmt"
	"net"
)

type HumanLogger struct{}

func (*HumanLogger) Accept(c net.Conn) bool {
	fmt.Printf("=> %v => %v\n", c.RemoteAddr(), c.LocalAddr())
	return true
}

func (*HumanLogger) Write(p []byte) error {
	fmt.Printf("-> %q\n", p)
	return nil
}

func (*HumanLogger) Read(p []byte) error {
	fmt.Printf("<- %q\n", p)
	return nil
}

func (*HumanLogger) String() string {
	return "human-logger"
}

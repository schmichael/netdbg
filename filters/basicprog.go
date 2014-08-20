package filters

import (
	"fmt"
	"io"
	"net"
)

func init() {
	RegisterFilter(basicProgressName, newBasicProgress)
}

var (
	basicProgSymbol = map[Role]string{
		Client: "→",
		Server: "←",
	}
)

const basicProgressName = "prog"

type BasicProgress struct {
	in  <-chan []byte
	out chan<- []byte
}

func newBasicProgress(r Role, i <-chan []byte, o chan<- []byte) Filter {
	p := &BasicProgress{in: i, out: o}
	go p.handle(basicProgSymbol[r])
	return p
}

func (*BasicProgress) Accept(net.Conn) bool {
	fmt.Print("⇄")
	return true
}

func (p *BasicProgress) handle(sym string) {
	var b []byte
	var ok bool
	for {
		if b, ok = <-p.in; !ok {
			close(p.out)
			return
		}
		fmt.Print(sym)
		p.out <- b
	}
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
	return basicProgressName
}

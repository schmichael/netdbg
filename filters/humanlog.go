package filters

import (
	"fmt"
	"io"
	"net"
	"time"
)

var (
	hlogSymbol = map[Role]string{
		Client: "⇨",
		Server: "⇦",
	}
)

func init() {
	RegisterFilter(humanLoggerName, newHumanLogger)
}

const humanLoggerName = "log"

type HumanLogger struct {
	in  <-chan []byte
	out chan<- []byte

	role  Role
	start time.Time
	count uint64
}

func newHumanLogger(r Role, i <-chan []byte, o chan<- []byte) Filter {
	h := HumanLogger{
		in:    i,
		out:   o,
		role:  r,
		start: time.Now(),
	}
	go h.handle(hlogSymbol[r])
	return &h
}

func (h *HumanLogger) Accept(c net.Conn) bool {
	fmt.Printf("  %v ⇄ %v\n", c.RemoteAddr(), c.LocalAddr())
	return true
}

func (h *HumanLogger) handle(sym string) {
	for {
		p, ok := <-h.in
		if !ok {
			close(h.out)
			return
		}

		h.count += uint64(len(p))
		fmt.Printf("%s %q\n", sym, p)
		h.out <- p
	}
}

func (h *HumanLogger) Close(err error) bool {
	dur := time.Now().Sub(h.start)
	var msg string
	if err == io.EOF {
		msg = "↹ closed"
	} else {
		msg = fmt.Sprintf("↯ %v", err)
	}
	fmt.Printf("%s: %s after %s; bytes: %d\n", h.role, msg, dur, h.count)
	return true
}

func (*HumanLogger) String() string {
	return "human-logger"
}

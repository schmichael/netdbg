package filters

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/schmichael/netdbg"
)

func init() {
	RegisterFilter(humanLoggerName, newHumanLogger)
}

const humanLoggerName = "log"

type HumanLogger struct {
	start time.Time
	sent  uint64
	recv  uint64
}

func newHumanLogger() netdbg.Filter { return &HumanLogger{} }

func (h *HumanLogger) Accept(c net.Conn) bool {
	h.start = time.Now()
	h.sent = 0
	h.recv = 0
	fmt.Printf("  %v ⇄ %v\n", c.RemoteAddr(), c.LocalAddr())
	return true
}

func (h *HumanLogger) Write(p []byte) error {
	h.sent += uint64(len(p))
	fmt.Printf("⇒ %q\n", p)
	return nil
}

func (h *HumanLogger) Read(p []byte) error {
	h.recv += uint64(len(p))
	fmt.Printf("⇒ %q\n", p)
	return nil
}

func (h *HumanLogger) Close(err error) bool {
	dur := time.Now().Sub(h.start)
	var msg string
	if err == io.EOF {
		msg = "↹ closed"
	} else {
		msg = fmt.Sprintf("↯ %v", err)
	}
	fmt.Printf("%s after %s; sent: %d  recv: %d\n", msg, dur, h.sent, h.recv)
	return true
}

func (*HumanLogger) String() string {
	return "human-logger"
}

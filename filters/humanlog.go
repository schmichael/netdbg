package filters

import (
	"fmt"
	"io"
	"net"
	"time"
)

func init() {
	RegisterFilter(humanLoggerName, newHumanLogger)
}

const humanLoggerName = "log"

type HumanLogger struct {
	writerIn  chan []byte
	writerOut chan []byte
	readerIn  chan []byte
	readerOut chan []byte

	start time.Time
	sent  uint64
	recv  uint64
}

func newHumanLogger(win chan []byte, rin chan []byte) (Filter, chan []byte, chan []byte) {
	h := HumanLogger{}
	h.writerIn = win
	h.writerOut = make(chan []byte)
	h.readerIn = rin
	h.readerOut = make(chan []byte)
	go h.write()
	go h.read()
	return &h, h.writerOut, h.readerOut
}

func (h *HumanLogger) Accept(c net.Conn) bool {
	h.start = time.Now()
	h.sent = 0
	h.recv = 0
	fmt.Printf("  %v ⇄ %v\n", c.RemoteAddr(), c.LocalAddr())
	return true
}

func (h *HumanLogger) write() {
	for {
		p, ok := <-h.writerIn
		if !ok {
			close(h.writerOut)
			return
		}

		h.sent += uint64(len(p))
		fmt.Printf("⇨ %q\n", p)
		h.writerOut <- p
	}
}

func (h *HumanLogger) read() {
	for {
		p, ok := <-h.readerIn
		if !ok {
			close(h.readerOut)
			return
		}

		h.recv += uint64(len(p))
		fmt.Printf("⇦ %q\n", p)
		h.readerOut <- p
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
	fmt.Printf("%s after %s; sent: %d  recv: %d\n", msg, dur, h.sent, h.recv)
	return true
}

func (*HumanLogger) String() string {
	return "human-logger"
}

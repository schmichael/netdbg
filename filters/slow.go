package filters

import (
	"bytes"
	"net"
	"time"
)

func init() {
	RegisterFilter(slowName, newSlowLoris)
}

const slowName = "slow"

type SlowLoris struct {
	in  <-chan []byte
	out chan<- []byte
}

func newSlowLoris(r Role, i <-chan []byte, o chan<- []byte) Filter {
	s := &SlowLoris{in: i, out: o}
	go s.handle()
	return s
}

func (s *SlowLoris) handle() {
	buf := bytes.NewBuffer(nil)
	for {
		// Wait until the buffer has data
		p, ok := <-s.in
		if !ok {
			close(s.out)
			return
		}
		buf.Write(p)

		// Slowly empty the buffer until no data is left
		for buf.Len() > 0 {
			select {
			case p, ok := <-s.in:
				if !ok {
					close(s.out)
					return
				}
				buf.Write(p)
			case <-time.After(time.Second):
				// Slowly write
				b, _ := buf.ReadByte()
				s.out <- []byte{b}
			}
		}
	}
}

func (*SlowLoris) Accept(net.Conn) bool {
	return true
}

func (*SlowLoris) Close(error) bool {
	return true
}

func (*SlowLoris) String() string {
	return slowName
}

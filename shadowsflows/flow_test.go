package shadowsflows

import (
	"net"
	"time"
	"testing"
)

type conn struct {
	r int
	w int
}

func (c *conn) Read(_ []byte) (n int, err error)   { return c.r, nil }
func (c *conn) Write(_ []byte) (n int, err error)  { return c.w, nil }
func (c *conn) Close() error                       { return nil }
func (c *conn) LocalAddr() (net.Addr)              { return nil }
func (c *conn) RemoteAddr() (net.Addr)             { return nil }
func (c *conn) SetDeadline(_ time.Time) error      { return nil }
func (c *conn) SetReadDeadline(_ time.Time) error  { return nil }
func (c *conn) SetWriteDeadline(_ time.Time) error { return nil }

func TestFlow_OnTrafficLimit(t *testing.T) {
	f := New()
	f.interval = time.Millisecond * 100 //0.1s
	pass := false
	f.OnTrafficLimit(func() {
		pass = true
	})
	f.SetTrafficLimit(1024 * 10)
	c := f.ReplaceConn(func(i net.Conn) net.Conn {
		return i
	})(&conn{1024, 1024})
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			c.Read([]byte{})
		} else {
			c.Write([]byte{})
		}
	}
	time.Sleep(time.Millisecond * 150)
	if !pass {
		t.FailNow()
	}
}

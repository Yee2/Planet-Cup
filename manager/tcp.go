package manager

import (
	"io"
	"net"
	"time"

	"encoding/binary"
	"github.com/shadowsocks/go-shadowsocks2/socks"
)

// Listen on addr for incoming connections.
func tcpRemote(addr string, shadow func(net.Conn) net.Conn) (io.Closer, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Danger("failed to listen on %s: %v", addr, err)
		return nil, err
	}
	closed := false
	stop := make(ch)

	logger.Success("listening TCP on %s", addr)

	go func() {
		<-stop
		closed = true
		l.Close()
	}()

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				if closed {
					stop <- 1
					break
				}
				logger.Danger("failed to accept: %v", err)
				continue
			}

			go func() {
				defer c.Close()
				c.(*net.TCPConn).SetKeepAlive(true)
				c = shadow(c)

				tgt, err := socks.ReadAddr(c)
				if err != nil {
					logger.Danger("failed to get target address: %v", err)
					return
				}
				//拦截访问面板地址
				target := "planet.cup"

				h := string(tgt[2 : 2+int(tgt[1])])
				p := binary.BigEndian.Uint16(tgt[2+int(tgt[1]) : 4+int(tgt[1])])
				if tgt[0] == 0x03 && int(tgt[1]) == len(target) && h == target && p == 80 {
					tgt = []byte{0x01, 0x7f, 0x00, 0x00, 0x01, 0x00, 0x00}
					binary.BigEndian.PutUint16(tgt[len(tgt)-2:len(tgt)], 34567)
				}
				rc, err := net.Dial("tcp", tgt.String())
				if err != nil {
					logger.Danger("failed to connect to target: %v", err)
					return
				}
				defer rc.Close()
				rc.(*net.TCPConn).SetKeepAlive(true)

				logger.Info("proxy %s <-> %s", c.RemoteAddr(), tgt)
				_, _, err = relay(c, rc)
				if err != nil {
					if err, ok := err.(net.Error); ok && err.Timeout() {
						return // ignore i/o timeout
					}
					logger.Danger("relay error: %v", err)
				}
			}()
		}
	}()
	return stop, nil
}

// relay copies between left and right bidirectionally. Returns number of
// bytes copied from right to left, from left to right, and any error occurred.
func relay(left, right net.Conn) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)

	go func() {
		n, err := io.Copy(right, left)
		right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
		left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
		ch <- res{n, err}
	}()

	n, err := io.Copy(left, right)
	right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
	left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
	rs := <-ch

	if err == nil {
		err = rs.Err
	}
	return n, rs.N, err
}

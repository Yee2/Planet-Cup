package shadowsflows

import (
	"net"
)


type Flows struct {
	net.Conn
	flow *Flow
}

func (self * Flows)Read(b []byte) (n int, err error){
	n, err = self.Conn.Read(b)
	self.flow.Up += n
	return n,err
}

func (self * Flows)Write(b []byte) (n int, err error)  {
	n, err = self.Conn.Write(b)
	self.flow.Down += n
	return n,err
}

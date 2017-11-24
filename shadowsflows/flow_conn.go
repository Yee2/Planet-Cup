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
	self.flow.Mu.Lock()
	self.flow.Up += n
	self.flow.Mu.Unlock()
	return n,err
}

func (self * Flows)Write(b []byte) (n int, err error)  {
	n, err = self.Conn.Write(b)
	self.flow.Mu.Lock()
	self.flow.Down += n
	self.flow.Mu.Unlock()
	return n,err
}

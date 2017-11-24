package shadowsflows

import (
	"net"
)

type Flows_PacketConn struct {
	net.PacketConn
	flow *Flow
}


func (self *Flows_PacketConn)ReadFrom(b []byte) (n int, addr net.Addr, err error){
	n,addr,err = self.PacketConn.ReadFrom(b)
	self.flow.Mu.Lock()
	self.flow.Up += n
	self.flow.Mu.Unlock()
	return
}
func (self *Flows_PacketConn)WriteTo(b []byte, addr net.Addr) (n int, err error){
	n,err = self.PacketConn.WriteTo(b,addr)
	self.flow.Mu.Lock()
	self.flow.Down += n
	self.flow.Mu.Unlock()
	return
}
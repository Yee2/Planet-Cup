package shadowsflows

import (
	"fmt"
	"code.cloudfoundry.org/bytefmt"
	"net"
)


type Flow struct {
	Up int
	Down int
}

func (self Flow)String() string {
	return fmt.Sprintf("上传流量:%s,下载流量:%s\n",bytefmt.ByteSize(uint64(self.Up)),bytefmt.ByteSize(uint64(self.Down)))
}
func New() *Flow {
	return &Flow{}
}
func (self *Flow)PipePacket(C net.PacketConn) net.PacketConn{
	return &Flows_PacketConn{C,self}
}

func (self *Flow)Pipe(C net.Conn) net.Conn{
	return &Flows{C,self}
}

func (self *Flow)ReplaceConn(shadow func(net.Conn) net.Conn)(func(net.Conn) net.Conn){
	return func(C net.Conn) net.Conn{
		return shadow(self.Pipe(C))
	}
}
func (self *Flow)ReplacePacketConn(shadow func(net.PacketConn) net.PacketConn)(func(net.PacketConn) net.PacketConn){
	return func(C net.PacketConn) net.PacketConn{
		return shadow(self.PipePacket(C))
	}
}

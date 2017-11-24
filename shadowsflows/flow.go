package shadowsflows

import (
	"fmt"
	"code.cloudfoundry.org/bytefmt"
	"net"
	"sync"
)


type Flow struct {
	Up int
	Down int
	Mu sync.Mutex
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

// 清空流量,并返回
func (self *Flow)Get()(Up int,Down int){
	self.Mu.Lock()
	Up,Down,self.Up,self.Down = self.Up,self.Down,0,0
	self.Mu.Unlock()
	return
}

// 原先流量加上新设置流量
func (self *Flow)Set(Up int,Down int){
	self.Mu.Lock()
	self.Up += Up
	self.Down += Down
	self.Mu.Unlock()
}
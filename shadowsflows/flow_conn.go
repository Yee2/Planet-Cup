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
	async(self.flow.speed.tcp.up,n)
	return n,err
}

// 将数据写入到用户端，对用户来说，就是下载流量，对服务器来说就是上传流量
func (self * Flows)Write(b []byte) (n int, err error)  {
	n, err = self.Conn.Write(b)
	self.flow.Mu.Lock()
	self.flow.Down += n
	self.flow.Mu.Unlock()
	async(self.flow.speed.tcp.down,n)
	return n,err
}

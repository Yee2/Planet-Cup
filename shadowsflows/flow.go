package shadowsflows

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"net"
	"sync"
	"time"
)

type Flow struct {
	Up             int
	Down           int
	Mu             sync.Mutex
	speed          speed
	BandwidthLimit int //速度限制，使用int类型，可以在32位支持2T左右最大值，64位无需担心
	TrafficLimit   int //总流量限制
	onTrafficLimit func()
	interval       time.Duration
	monitor        bool
}

func (self Flow) String() string {
	return fmt.Sprintf("上传流量:%s,下载流量:%s\n", bytefmt.ByteSize(uint64(self.Up)), bytefmt.ByteSize(uint64(self.Down)))
}

func New() *Flow {
	self := &Flow{interval: time.Minute * 10}
	self.speed.run()
	return self
}
func (self *Flow) OnTrafficLimit(callback func()) {
	if callback != nil {
		self.Mu.Lock()
		self.onTrafficLimit = callback
		self.Mu.Unlock()
	}
}
func (self *Flow) SetTrafficLimit(i int) {
	self.Mu.Lock()
	defer self.Mu.Unlock()
	self.TrafficLimit = i
	if i > 0 && self.monitor == false {
		self.Deamon()
	}
}
func (self *Flow) Deamon() {
	go func() {
		self.monitor = true
		defer func() { self.monitor = false }()
		for {
			if self.TrafficLimit <= 0 {
				return
			}
			if self.TrafficLimit <= self.Up+self.Down {
				if self.onTrafficLimit != nil {
					self.onTrafficLimit()
				}
			}
			time.Sleep(self.interval)
		}
	}()
}
func (self *Flow) PipePacket(C net.PacketConn) net.PacketConn {
	return &Flows_PacketConn{C, self}
}
func (self *Flow) Speed() (Up int, Down int) {
	return self.speed.tcp.up.bucket + self.speed.udp.up.bucket,
		self.speed.tcp.down.bucket + self.speed.udp.down.bucket
}
func (self *Flow) SpeedHistory(n int) []struct {
	Up   int
	Down int
	Time int64
} {
	h := self.speed.history
	t := len(h.line)
	var i int
	if n > t {
		n = t // 不能超出数组
	}
	n--
	data := make([]struct {
		Up   int
		Down int
		Time int64
	}, 0, n)
	// 从指针开始处往前移动 n -1个
	for n = h.offset - n; n <= h.offset; n++ {

		if n < 0 {
			i = t + n
		} else {
			i = n
		}

		item := h.line[i]
		if item.time.Unix() < 0 {
			continue
		}
		data = append(data, struct {
			Up   int
			Down int
			Time int64
		}{item.up, item.down, item.time.Unix()})
	}
	return data
}
func (self *Flow) Pipe(C net.Conn) net.Conn {
	return &Flows{C, self}
}

func (self *Flow) ReplaceConn(shadow func(net.Conn) net.Conn) func(net.Conn) net.Conn {
	return func(C net.Conn) net.Conn {
		return shadow(self.Pipe(C))
	}
}
func (self *Flow) ReplacePacketConn(shadow func(net.PacketConn) net.PacketConn) func(net.PacketConn) net.PacketConn {
	return func(C net.PacketConn) net.PacketConn {
		return shadow(self.PipePacket(C))
	}
}

// 清空流量,并返回
func (self *Flow) Get() (Up int, Down int) {
	self.Mu.Lock()
	Up, Down, self.Up, self.Down = self.Up, self.Down, 0, 0
	self.Mu.Unlock()
	return
}

// 原先流量加上新设置流量
func (self *Flow) Set(Up int, Down int) {
	self.Mu.Lock()
	self.Up += Up
	self.Down += Down
	self.Mu.Unlock()
}

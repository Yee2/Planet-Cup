package shadowsflows

import (
	"fmt"
	"code.cloudfoundry.org/bytefmt"
	"net"
	"sync"
	"time"
)

type Flow struct {
	Up    int
	Down  int
	Mu    sync.Mutex
	speed speed
}

// 神经病一样的写法
type speed struct {
	tcp struct {
		up   addup
		down addup
	}
	udp struct {
		up   addup
		down addup
	}
	history struct {
		line [600]struct {
			up   int
			down int
			time time.Time
		} //记录十分钟内的数据
		offset int
		sync.Mutex
	}
}

func (self *speed) run() {
	self.udp.up.ch = make(chan int, 100)
	self.udp.down.ch = make(chan int, 100)
	self.tcp.up.ch = make(chan int, 100)
	self.tcp.down.ch = make(chan int, 100)
	self.udp.up.run()
	self.udp.down.run()
	self.tcp.up.run()
	self.tcp.down.run()
	go func() {
		for {
			<-time.Tick(time.Second)
			self.history.Lock()
			self.history.line[self.history.offset].up = self.tcp.up.bucket + self.udp.up.bucket
			self.history.line[self.history.offset].down = self.tcp.down.bucket + self.udp.down.bucket
			self.history.line[self.history.offset].time = time.Now()
			self.history.offset ++
			if self.history.offset >= len(self.history.line){
				self.history.offset -= len(self.history.line)
			}
			self.history.Unlock()
		}
	}()
}

type addup struct {
	current int
	bucket  int
	sync.Mutex
	ch      chan int
}

func (self Flow) String() string {
	return fmt.Sprintf("上传流量:%s,下载流量:%s\n", bytefmt.ByteSize(uint64(self.Up)), bytefmt.ByteSize(uint64(self.Down)))
}

func (self *addup) run() {
	go func() {
		for {
			<-time.Tick(time.Second)
			self.Lock()
			self.bucket = self.current
			self.current = 0
			self.Unlock()

		}
	}()
	go func() {
		for {
			n := <-self.ch
			self.Lock()
			self.current += n
			self.Unlock()
		}
	}()
}
func New() *Flow {
	self := &Flow{}
	self.speed.run()
	return self
}
func (self *Flow) PipePacket(C net.PacketConn) net.PacketConn {
	return &Flows_PacketConn{C, self}
}
func (self *Flow) Speed() (Up int, Down int) {
	return self.speed.tcp.up.bucket + self.speed.udp.up.bucket,
		self.speed.tcp.down.bucket + self.speed.udp.down.bucket
}
func (self *Flow) Speed_history(n int) ([]struct{Up int;Down int;Time int64}) {
	h := self.speed.history
	t := len(h.line)
	var i int
	if n > t{
		n = t// 不能超出数组
	}
	n--
	data := make([]struct{Up int;Down int;Time int64},0,n)
	// 从指针开始处往前移动 n -1个
	for n = h.offset - n; n <= h.offset; n++{

		if n < 0 {
			i = t + n
		}else{
			i = n
		}

		item := h.line[i]
		if item.time.Unix() < 0{
			continue
		}
		data = append(data,struct{Up int;Down int;Time int64}{item.up,item.down,item.time.Unix()})
	}
	return data
}
func (self *Flow) Pipe(C net.Conn) net.Conn {
	return &Flows{C, self}
}

func (self *Flow) ReplaceConn(shadow func(net.Conn) net.Conn) (func(net.Conn) net.Conn) {
	return func(C net.Conn) net.Conn {
		return shadow(self.Pipe(C))
	}
}
func (self *Flow) ReplacePacketConn(shadow func(net.PacketConn) net.PacketConn) (func(net.PacketConn) net.PacketConn) {
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

// 通过通道，刷新网速，同时不至于卡死
func async(speed addup, i int) {
	go func() {
		select {
		case speed.ch <- i:
		case <-time.After(time.Second):
		}
	}()
}

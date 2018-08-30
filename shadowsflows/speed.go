package shadowsflows

import (
	"time"
	"sync"
)

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
			time.Sleep(time.Second)
			self.history.Lock()
			self.history.line[self.history.offset].up = self.tcp.up.bucket + self.udp.up.bucket
			self.history.line[self.history.offset].down = self.tcp.down.bucket + self.udp.down.bucket
			self.history.line[self.history.offset].time = time.Now()
			self.history.offset ++
			if self.history.offset >= len(self.history.line) {
				self.history.offset -= len(self.history.line)
			}
			self.history.Unlock()
		}
	}()
}

// 这是一个计数器
type addup struct {
	current int      // 前一秒走过的流量
	bucket  int      // 当前秒走的流量
	ch      chan int // 没过一秒刷新一次
	sync.Mutex
}

func (self *addup) run() {
	go func() {
		for {
			time.Sleep(time.Second)
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

// 通过通道，刷新网速，同时不至于卡死
func async(speed addup, i int) {
	go func() {
		select {
		case speed.ch <- i:
		case <-time.After(time.Second):
		}
	}()
}

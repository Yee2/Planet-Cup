package ylog

import (
	"errors"
	"sync"
	"time"
)

var (
	ReceiverClosed = errors.New("ReceiverClosed")
)

type Receiver struct {
	ch     chan *Log
	closed chan int
	dead   bool
}

func (self *Receiver) Receive() (log *Log, err error) {
	select {
	case log = <-self.ch:
		return log, nil
	case <-self.closed:
		self.Close()
		return nil, ReceiverClosed
	}
}
func (self *Receiver) Close() {
	receivers.Lock()
	defer receivers.Unlock()
	for i := range receivers.list {
		if receivers.list[i] == self {
			receivers.list = append(receivers.list[:i], receivers.list[i:]...)
		}
	}
	if self.dead == false {
		close(self.ch)
		close(self.closed)
	}
	self.dead = true
}
func NewReceiver() *Receiver {
	r := &Receiver{make(chan *Log), make(chan int), false}
	receivers.Lock()
	defer receivers.Unlock()
	receivers.list = append(receivers.list, r)
	return r
}

var receivers = struct {
	list []*Receiver
	sync.Mutex
}{make([]*Receiver, 0), sync.Mutex{}}

func broadcast(log *Log) {
	for _, v := range receivers.list {
		go func(r *Receiver) {
			select {
			case r.ch <- log:
			case <-time.After(time.Second * 2):
			}
		}(v)
	}
}

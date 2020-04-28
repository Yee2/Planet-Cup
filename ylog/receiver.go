package ylog

import (
	"errors"
	"sync"
)

var (
	ReceiverClosed = errors.New("ReceiverClosed")
)

type Receiver struct {
	ch     chan *Log
	closed chan int
	dead   bool
}

func (r *Receiver) Receive() (log *Log, err error) {
	select {
	case log = <-r.ch:
		return log, nil
	case <-r.closed:
		r.Close()
		return nil, ReceiverClosed
	}
}
func (r *Receiver) Close() {
	receivers.Lock()
	defer receivers.Unlock()
	for i := range receivers.list {
		if receivers.list[i] == r {
			receivers.list = append(receivers.list[:i], receivers.list[i:]...)
		}
	}
	if r.dead == false {
		close(r.ch)
		close(r.closed)
	}
	r.dead = true
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


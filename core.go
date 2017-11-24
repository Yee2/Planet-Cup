package main

import (
	"github.com/Yee2/ssMulti-user/shadowsflows"
	"log"
	"errors"
	"github.com/riobard/go-shadowsocks2/core"
	"fmt"
)

type Shadowsocks struct {
	Addr string
	Port int
	Password string
	Cipher string
	shadowsflows.Flow
}
func NewTable() *Table{
	return &Table{make(map[int]chan int),make(map[int]*Shadowsocks),make([]byte,0,32),""}
}
type Table struct {
	chans map[int]chan int
	rows map[int]*Shadowsocks
	key []byte
	remote string
}

func (self *Table)start(id int)error {
	if ss,ok := self.rows[id];ok{
		if _,ok := self.chans[id]; ok{
			return nil
		}
		ciph, err := core.PickCipher(ss.Cipher, self.key, ss.Password)
		if err != nil {
			return err
		}
		self.chans[id] = make(chan int)
		go udpRemote(self.chans[id],ss.Addr, ss.ReplacePacketConn(ciph.PacketConn))
		go tcpRemote(self.chans[id],ss.Addr, ss.ReplaceConn(ciph.StreamConn))
		return nil
	}
	return errors.New(fmt.Sprintf("%d Not Exist!",id))
}
func (self *Table)stop(id int){
	logf("Stop:%d",id)
	if c,ok := self.chans[id];ok{
		c <- 1
		c <- 1
	}
	delete(self.chans,id)
}
func (self *Table)boot(){
	for id,_ := range self.rows{
		logf("Listen on %d\n",id)
		err := self.start(id)
		if err != nil{
			log.Fatal(err)
		}
	}
}
func (self *Table)push(ss *Shadowsocks) error {
	if ss.Port < 1{
		return errors.New("Port error!")
	}
	if _,ok := self.rows[ss.Port];ok{
		return errors.New("Port Existed!")
	}
	self.rows[ss.Port] = ss
	return nil
}

func (self *Table)add(ss *Shadowsocks) error {
	if err := self.push(ss);err != nil{
		return err
	}
	if err := self.start(ss.Port);err != nil{
		return err
	}
	return nil
}
func (self *Table)pwd(id int,password string) error {
	ss,ok := self.rows[id]
	if ok{
		return errors.New("Port Existed!")
	}
	self.stop(id)
	ss.Password = password
	return self.start(id)
}
func (self *Table)del(id int) error {
	self.stop(id)
	if _,ok := self.rows[id]; ok{
		delete(self.rows,id)
		return nil
	}
	return errors.New(fmt.Sprintf("Port %d No Found!",id))
}



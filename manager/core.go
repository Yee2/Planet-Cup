// 提供快捷的shadowsocks服务管理接口
package manager

import (
	"errors"
	"fmt"

	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/Yee2/Planet-Cup/shadowsflows"
	"net"
	"strconv"
	"io"
)

type Shadowsocks struct {
	Addr     string
	Port     int
	Password string
	Cipher   string
	*shadowsflows.Flow
	key      []byte
	tcp      io.Closer
	udp      io.Closer
}

func NewShadowsocks(port int, password, method string) (*Shadowsocks, error) {
	if port < 1 || port > 65535 {
		return nil, errors.New("Port error!")
	}
	return &Shadowsocks{
		Port:     port,
		Password: password,
		Cipher:   method,
		Flow:     shadowsflows.New(),
		key:      make([]byte, 0, 32),
	}, nil
}

func (self *Shadowsocks) Start() error {
	if self.udp != nil || self.tcp != nil {
		return nil
	}
	ciph, err := core.PickCipher(self.Cipher, self.key, self.Password)
	if err != nil {
		return err
	}

	if closer, err := udpRemote(net.JoinHostPort(self.Addr, strconv.Itoa(self.Port)), self.ReplacePacketConn(ciph.PacketConn)); err != nil {
		return err
	} else {
		self.udp = closer
	}
	if closer, err := tcpRemote(net.JoinHostPort(self.Addr, strconv.Itoa(self.Port)), self.ReplaceConn(ciph.StreamConn)); err != nil {
		return err
	} else {
		self.tcp = closer
	}
	return nil
}

func (self *Shadowsocks) Stop() error {
	if self.tcp != nil {
		if err := self.tcp.Close(); err != nil {
			logger.Warning("停止TCP监听发生错误:%s，端口:%d",err,self.Port)
			return err
		}
		self.tcp = nil
	}
	if self.udp != nil {
		if err := self.udp.Close(); err != nil {
			logger.Warning("停止UDP监听发生错误:%s，端口:%d",err,self.Port)
			return err
		}
		self.udp = nil
	}
	return nil
}
func (self *Shadowsocks) String() string {
	return fmt.Sprintf("ss://%s:%s@%s:%d", self.Cipher, self.Password, self.Addr, self.Port)
}

type Table struct {
	Rows   map[int]*Shadowsocks
	key    []byte
	remote string
}

func NewTable() *Table {
	return &Table{make(map[int]*Shadowsocks), make([]byte, 0, 32), ""}
}

// 这个方法将被弃用
func (self *Table) Start(id int) error {
	if ss, ok := self.Rows[id]; ok {
		return ss.Start()
	}
	return errors.New(fmt.Sprintf("%d Not Exist!", id))
}

// 这个方法将被弃用
func (self *Table) Stop(id int) {
	if ss, ok := self.Rows[id]; ok {
		ss.Stop()
	}
}

// 启动全部服务
func (self *Table) Boot() {
	for _, ss := range self.Rows {
		err := ss.Start()
		if err != nil {
			logger.Warning("启动服务发生错误:%s", err)
		}
	}
}

// 关闭全部服务
func (self *Table) Shutdown() {
	for _, ss := range self.Rows {
		ss.Stop()
	}
}

// 添加一个新的SS服务
func (self *Table) Push(ss *Shadowsocks) error {
	if _, ok := self.Rows[ss.Port]; ok {
		return errors.New("Port Existed!")
	}
	self.Rows[ss.Port] = ss
	return nil
}

// 添加一个新的服务并启动这个服务
func (self *Table) Add(ss *Shadowsocks) error {

	if _, ok := self.Rows[ss.Port]; ok {
		return errors.New("Port Existed!")
	}

	if err := ss.Start(); err != nil {
		return err
	}
	self.Rows[ss.Port] = ss
	return nil
}

// 修改一个SS服务的信息，端口必须保持一致
func (self *Table) Set(ss *Shadowsocks) error {
	if item, ok := self.Rows[ss.Port]; ok {
		if err := item.Stop(); err != nil {
			return err
		}
		self.Rows[ss.Port] = ss
		err := ss.Start()
		return err
	}
	return errors.New(fmt.Sprintf("Port(%d) does not exist!", ss.Port))
}
func (self *Table) Pwd(id int, password string) error {
	if ss, ok := self.Rows[id]; ok {
		if err := ss.Stop(); err != nil {
			return err
		}
		ss.Password = password
		if err := ss.Start(); err != nil {
			return err
		}
		return nil
	}
	return errors.New(fmt.Sprintf("Port(%d) does not exist!", id))
}

// 删除一个服务
func (self *Table) Del(id int) error {
	if ss, ok := self.Rows[id]; ok {
		ss.Stop()
		delete(self.Rows, id)
		return nil
	}
	return errors.New(fmt.Sprintf("Port %d No Found!", id))
}

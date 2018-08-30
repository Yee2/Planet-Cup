// 提供快捷的shadowsocks服务管理接口
package manager

import (
	"errors"
	"fmt"

	"encoding/json"
	"github.com/Yee2/Planet-Cup/shadowsflows"
	"github.com/shadowsocks/go-shadowsocks2/core"
	"io"
	"net"
	"os"
	"strconv"
	"math"
)

const save_flag = "YeeShadowSocksPanel_v0.1"

type Shadowsocks struct {
	Addr     string
	Port     int
	Password string
	Cipher   string
	key      []byte
	tcp      io.Closer
	udp      io.Closer
	*shadowsflows.Flow
	status   Status
}
type Status int

const (
	Active   Status = iota //运行中
	Suspend                //超过限制被暂停
	Shutdown               // 未运行
	Failed                 //启动失败
	Unknown                // 未知状态，关闭失败的时候返回
)

func (s Status)String() string {
	switch s {
	case Active:
		return "运行中"
	case Shutdown:
		return "未运行"
	case Unknown:
		return "未知错误"
	case Suspend:
		return "超过限额"
	case Failed:
		return "启动失败"
	default:
		return fmt.Sprintf("Unknown(%02X)",s)
	}
}
func NewShadowsocks(port int, password, method string) (*Shadowsocks, error) {
	if port < 1 || port > 0xffff {
		return nil, errors.New("port error")
	}
	ss := &Shadowsocks{
		Port:     port,
		Password: password,
		Cipher:   method,
		Flow:     shadowsflows.New(),
		key:      make([]byte, 0, 32),
		status:   Shutdown,
	}
	ss.OnTrafficLimit(func() {
		ss.tcp.Close()
		ss.udp.Close()
		ss.status = Suspend
	})
	return ss, nil
}
func (self *Shadowsocks) SetTrafficLimit(i int){
	self.Flow.SetTrafficLimit(i*int(math.Pow(2,30)))
}
func (self *Shadowsocks) GetTrafficLimit()(i int){
	return self.Flow.TrafficLimit/int(math.Pow(2,30))
}
func (self *Shadowsocks) Status() Status {
	return self.status
}
func (self *Shadowsocks) Start() (e error) {
	defer func() {
		if e != nil {
			self.status = Failed
		} else {
			self.status = Active
		}
	}()
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

func (self *Shadowsocks) Stop() (e error) {
	defer func() {
		if e != nil {
			self.status = Unknown
		} else {
			self.status = Shutdown
		}
	}()
	if self.tcp != nil {
		if err := self.tcp.Close(); err != nil {
			logger.Warning("停止TCP监听发生错误:%s，端口:%d", err, self.Port)
			return err
		}
		self.tcp = nil
	}
	if self.udp != nil {
		if err := self.udp.Close(); err != nil {
			logger.Warning("停止UDP监听发生错误:%s，端口:%d", err, self.Port)
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
	return fmt.Errorf("%d Not Exist!", id)
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

func (self *Table) Save(name string) (e error) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Truncate(0)
	//_, err = f.Write([]byte(save_flag))
	//if err != nil {
	//	return err
	//}
	//_, err = f.Write([]byte{0xff, 0xff})
	//if err != nil {
	//	return err
	//}
	data := []struct {
		Port           int
		Password       string
		Cipher         string
		Up             int
		Down           int
		TrafficLimit int
	}{}
	for _, ss := range self.Rows {
		data = append(data, struct {
			Port           int
			Password       string
			Cipher         string
			Up             int
			Down           int
			TrafficLimit int
		}{ss.Port, ss.Password, ss.Cipher, ss.Flow.Up, ss.Flow.Down, ss.TrafficLimit})
	}
	str, err := json.Marshal(data)
	_, err = f.Write(str)
	return err
}

func (self *Table) Load(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	data := []struct {
		Port           int
		Password       string
		Cipher         string
		Up             int
		Down           int
		TrafficLimit int
	}{}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	for i := range data {
		ss, err := NewShadowsocks(data[i].Port, data[i].Password, data[i].Cipher)
		if err != nil {
			logger.Danger("%s", err)
		}
		ss.Flow.Up = data[i].Up
		ss.Flow.Down = data[i].Down
		ss.SetTrafficLimit(data[i].TrafficLimit)
		self.Add(ss)
	}
	return nil
}

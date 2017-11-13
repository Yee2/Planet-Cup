package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yee2/ssMulti-user/shadowsflows"
	"strings"
	"strconv"
	"fmt"
)

var config struct {
	Verbose    bool
	UDPTimeout time.Duration
}
var shadowsTable = NewTable()
func logf(f string, v ...interface{}) {
	if config.Verbose {
		log.Printf(f, v...)
	}
}

func main() {

	log.SetFlags( log.Ldate | log.Ltime | log.Llongfile )
	args := os.Args[1:]
	for _,u := range args{

		addr, cipher, password , err := parseURL(u)
		if err!=nil{
			log.Fatal(err)
		}
		p,err := GetPort(addr)
		if err!=nil{
			log.Fatal(err)
		}
		ss := &Shadowsocks{Addr:addr,Port:p,Password:password,Cipher:cipher,Flow:shadowsflows.Flow{0,0}}
		shadowsTable.push(ss)
	}
	shadowsTable.boot()
	go shadowsTable.Listen()
	go func() {
		for  {
			// TODO: 好好学习终端控制符
			n := len(shadowsTable.rows)
			fmt.Print("\x1b[2J\x1b[1;1H===================\n")
			for id,ss := range shadowsTable.rows{
				fmt.Printf("\x1b[K[%d]%s",id,ss.Flow)
			}
			fmt.Print("\x1b[K===================\n")
			//向上移动光标
			fmt.Printf("\x1b[%dA",n+2)
			//fmt.Printf("\x1b[%d;H",n)
			time.Sleep(3 * time.Second)
		}
	}()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func parseURL(s string) (addr, cipher, password string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}

	addr = u.Host
	if u.User != nil {
		cipher = u.User.Username()
		password, _ = u.User.Password()
	}
	return
}

func GetPort(s string)(int,error){

	str := strings.Split(s,":")
	p,err := strconv.ParseInt(str[len(str) - 1],10,64)
	if err != nil{
		return 0,err
	}
	return int(p),nil
}
// 让程序作为一个服务端来运行
package server

import (
	"encoding/json"
	"fmt"
	"github.com/Yee2/Planet-Cup/shadowsflows"
	"log"
	"net"
)

func (self *Table) Listen() {
	ln, err := net.Listen("tcp", ":8088")
	logf("Start listening on port 8088")
	if err != nil {
		// handle error
		//logf(err.Error())
		log.Fatalf("failed to listen on %s: %v", ":8088", err.Error())
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			logf("%s", err)
			continue
		}
		go self.HandleConnection(conn)
	}
}

func (self *Table) HandleConnection(conn net.Conn) {
	var b []byte = make([]byte, 512)
	for {
		n, err := conn.Read(b)
		if err != nil {
			if err.(net.UnknownNetworkError).Timeout() {
				conn.Close()
				return
			}
			logf("%s", err)
			continue
		}
		cd := &Command{}
		res := make(map[string]interface{})
		res_list := make([]interface{}, 0)
		err = json.Unmarshal(b[:n], cd)
		if err != nil {
			logf("%s", err)
		}

		switch cd.Command {
		case "add":
			ss := &Shadowsocks{fmt.Sprintf(":%d", cd.Port), cd.Port, cd.Password, cd.Cipher, shadowsflows.Flow{}}
			err = self.add(ss)
			if err != nil {
				res["error"] = err.Error()
			} else {
				res["port"] = cd.Port
				res["password"] = cd.Password
				res["cipher"] = cd.Cipher
			}
		case "list":
			for _, ss := range self.rows {
				res_list = append(res_list, struct {
					Port     int    `json:"port"`
					Password string `json:"password"`
					Cipher   string `json:"cipher"`
				}{ss.Port, ss.Password, ss.Cipher})
			}

		case "del":
			err = self.del(cd.Port)
			if err != nil {
				res["error"] = err.Error()
			} else {
				res["port"] = cd.Port
			}
		case "version":
			res["version"] = "0.0.1"
		case "pwd":
			err = self.pwd(cd.Port, cd.Password)
			if err != nil {
				res["error"] = err.Error()
			} else {
				res["port"] = cd.Port
				res["password"] = cd.Password
			}
		case "flow":
			for _, ss := range self.rows {
				res_list = append(res_list, struct {
					Port int `json:"port"`
					Up   int `json:"up"`
					Down int `json:"down"`
				}{ss.Port, ss.Flow.Up, ss.Flow.Down})
			}
		default:
			conn.Write([]byte("hello"))
			continue
		}
		txt := make([]byte, 0, 512)
		if len(res) > 0 {
			txt, err = json.Marshal(res)
		} else {
			txt, err = json.Marshal(res_list)
		}
		if err != nil {
			conn.Write([]byte("hello"))
			logf(err.Error())
		} else {
			conn.Write(txt)
		}
	}

}

type Command struct {
	Command  string                 `json:"command"`
	Port     int                    `json:"port"`
	Password string                 `json:"password"`
	Cipher   string                 `json:"cipher"`
	Options  map[string]interface{} `json:"options"`
}

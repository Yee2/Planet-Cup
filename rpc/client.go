package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Yee2/Planet-Cup/manager"
	"github.com/Yee2/Planet-Cup/shadowsflows"
	l "github.com/Yee2/Planet-Cup/ylog"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var Duration = time.Minute * 30
var n = 0
var logger = l.NewLogger("RPC")

type Client struct {
	url  string
	key  string
	ss   *manager.Table
	time time.Time /*最后更新时间时间*/
}

func NewClient(u, key string) (*Client, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", u+"/index.php/RESTful/test/"+key, nil)
	req.Header.Add("X-API-KEY", key)
	resp, err := client.Do(req)

	//resp,err := http.Get(u + "/index.php/RESTful/test/" + key)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("Authentication failed")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := struct {
		Result string `json:"result"`
		Error  string `json:"error"`
	}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	if res.Error != "" {
		return nil, errors.New(res.Error)
	}

	if res.Result == "done" {
		return &Client{u, key, manager.NewTable(), time.Now()}, nil
	}
	return nil, errors.New("未知错误!")
}

func (self *Client) q(p string, v interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", self.url+p, nil)
	req.Header.Add("X-API-KEY", self.key)
	resp, err := client.Do(req)

	//resp,err := http.Get(self.url + p)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	n++
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

func (self *Client) init() error {
	logger.Info("initialization ...")
	result := struct {
		Error string `json:"error"`
		List  []struct {
			Port     int    `json:"port"`
			Password string `json:"password"`
			Method   string `json:"method"`
		}
	}{}
	err := self.q("/index.php/RESTful/index", &result)
	if err != nil {
		return err
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	for _, item := range result.List {
		ss := &manager.Shadowsocks{fmt.Sprintf(":%d", item.Port), item.Port, item.Password,
			item.Method, shadowsflows.Flow{}}
		err := self.ss.Push(ss)
		if err != nil {
			logger.Info("%s", err)
		}
	}
	self.ss.Boot()
	logger.Info("initialization ...done")
	return nil
}

// 从服务器获取更新
func (self *Client) Pull() error {
	logger.Info("Pull data ...")
	result := struct {
		Error string `json:"error"`
		List  []struct {
			Key        string `json:"key"`
			Shadowsock struct {
				Port     int    `json:"port"`
				Password string `json:"password"`
				Method   string `json:"method"`
			} `json:"data"`
		}
	}{}
	err := self.q("/index.php/RESTful/action/"+fmt.Sprintf("%d", self.time.Unix()), &result)
	if err != nil {
		return err
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}

	for _, item := range result.List {
		switch item.Key {
		case "delete":
			logger.Info("Delete %d ...", item.Shadowsock.Port)
			self.ss.Del(item.Shadowsock.Port)
		case "put":
			logger.Info("set %d ...", item.Shadowsock.Port)
			ss := &manager.Shadowsocks{fmt.Sprintf(":%d", item.Shadowsock.Port), item.Shadowsock.Port,
				item.Shadowsock.Password, item.Shadowsock.Method, shadowsflows.Flow{}}
			err := self.ss.Set(ss)
			if err != nil {
				logger.Info("%s", err)
			}

		case "post":
			logger.Info("add %d ...", item.Shadowsock.Port)
			ss := &manager.Shadowsocks{fmt.Sprintf(":%d", item.Shadowsock.Port), item.Shadowsock.Port,
				item.Shadowsock.Password, item.Shadowsock.Method, shadowsflows.Flow{}}
			err := self.ss.Add(ss)
			if err != nil {
				logger.Info("%s", err)
			}
		default:
			logger.Info("Unknown command:%s", item.Key)
		}
	}
	self.time = time.Now()
	return nil
}

// 将本地流量推送到服务器上

func (self *Client) Push() (err error) {
	logger.Info("Push data ...")
	data := make([]struct {
		Port int `json:"port"`
		Up   int `json:"up"`
		Down int `json:"down"`
	}, 0, len(self.ss.Rows))
	for _, item := range self.ss.Rows {
		up, down := item.Get()
		data = append(data, struct {
			Port int `json:"port"`
			Up   int `json:"up"`
			Down int `json:"down"`
		}{item.Port, up, down})
	}
	defer func() {
		if e := recover(); e != nil {
			for _, item := range data {
				if flow, ok := self.ss.Rows[item.Port]; ok {
					flow.Set(item.Up, item.Down)
				}
			}
			err = e.(error)
		}
	}()
	postValues := url.Values{}
	str, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	postValues.Add("data", string(str))
	postValues.Add("time", fmt.Sprintf("%d", self.time.Unix()))
	req, err := http.NewRequest("POST", self.url+"/index.php/RESTful/push", strings.NewReader(postValues.Encode()))
	req.Header.Add("X-API-KEY", self.key)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 403 {
		panic(errors.New("Authentication failed!"))
	}
	defer resp.Body.Close()
	result := struct {
		Error   string `json:"error"`
		Result  string `json:"result"`
		Message string `json:"message"`
	}{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	if result.Error != "" {
		panic(errors.New(result.Error))
	}

	if result.Result != "done" {
		panic(errors.New(result.Message))
	}

	return nil
}

func (self *Client) Run() {
	if err := self.init(); err != nil {
		log.Fatal(err)
	}
	for {
		logger.Info("Sleep %s...", Duration.String())
		time.Sleep(Duration)
		if err := self.Push(); err != nil {
			logger.Info("%v", err)
		}

		if err := self.Pull(); err != nil {
			logger.Info("%v", err)
		}
	}
}

func (self *Client) Shutdown() {
	logger.Info("Shutdown ... ")
	self.Push()
	self.ss.Shutdown()
}

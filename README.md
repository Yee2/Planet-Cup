# 说明

基于`go-shadowsocks2`这个项目修改而来.


# 客户端模式

本程序作为客户端,将数据流量使用信息发送到服务器端,并从服务端获取数据,接口采用RESTful规范.

服务端返回标准json格式,有两个键值:检查`error`是否为空,不为空时候存在错误,没有错误的时候放回`list`为目标数据.
所有与shadowsocks相关的数据均采用一下结构:
````go
Shadowsock struct{
				Port int `json:"port"`
				Password string `json:"password"`
				Method string `json:"method"`
} 

````

# 以下内容已过期

# 食用方法

````shell
./ssMulti-user
````

# 接口协议

采用TCP协议,参考[shadowsocks-manager接口协议](https://github.com/shadowsocks/shadowsocks-manager/wiki/%E6%8E%A5%E5%8F%A3%E5%8D%8F%E8%AE%AE)
````go
type Command struct {
	Command string `json:"command"`
	Port int `json:"port"`
	Password string `json:"password"`
	Cipher string `json:"cipher"`
	Options map[string]interface{} `json:"options"`
}

````
## 注意

请使用双引号,键名需要使用双引号括起来,添加需要制定加密方式.

Command 实现的接口:`add` `del` `flow` `version` `pwd` `list`

## Bug

 ~~已知删除操作,需要客户端断开链接才能生效,UDP通道可能无法正常关闭.~~ 

## 返回错误

当发生错误时候,`error`可以查看错误信息.
````
{"error":"Port Existed!"}
````
## add 指令
````
{
  "command": "add",
  "port": 8388,
  "cipher": "AES-256-CFB",
  "password": "123456"
}
````
返回
````
{
  "port": 8388,
  "cipher": "AES-256-CFB",
  "password": "123456"
}
````

## del 指令
````
{
  "command": "del",
  "port": 8388
}
````
返回
````
{
  "port": 8388
}
````

## flow 指令
````
{
  "command: "flow"
}
````
返回：
````
[
  { "port": 1234, "up":1024,"down": 2048 },
  { "port": 1235, "up":1024,"down": 2048 }
]
````

## list 指令
````
{"command":"list"}
````
返回
````
[
{"port":6666,"password":"12345678","cipher":"AES-256-CFB"},
{"port":6667,"password":"12345678","cipher":"AES-256-CFB"}
]
````
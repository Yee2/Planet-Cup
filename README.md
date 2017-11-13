# 说明

基于`go-shadowsocks2`这个项目修改而来.

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

## add 指令
````
{
  "command": "add",
  "port": 8388,
  "cipher": "AES-256-CFB"
  "password": "123456"
}
返回
{
  "port": 8388,
  "cipher": "AES-256-CFB",
  "password": "123456"
}
````
## flow 指令
````
{
  "command: "flow"
}
返回：
[
  { port: 1234, up:1024,down: 2048 },
  { port: 1235, up:1024,down: 2048 }
]
````
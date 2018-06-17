# 说明

基于`go-shadowsocks2`这个项目修改而来，目的提供一个单用户简洁易用的管理面板。

## linux 下安装golang

适用于Debian和Redhat系列和其他小众Linux系统。
```sh
    #下载Golang安装包，版本 1.10
    wget https://dl.google.com/go/go1.10.2.linux-amd64.tar.gz
    # 解压安装包到官方推荐的目录
    tar -C /usr/local -xzf go1.10.2.linux-amd64.tar.gz
    
    # 修改环境变量，将 /usr/local/go/bin 添加进去
    vim /etc/environment
```

# 安装
```sh
    go get -u github.com/Yee2/Planet-Cup
```

# 使用

```sh
    # 由于路径问题，必须在下面目录执行命令
    cd ~/go/src/github.com/Yee2/Planet-Cup/
    ~/go/bin/Planet-Cup -c ~/go/src/github.com/Yee2/Planet-Cup/config.json
```
# 如何访问

面板默认并不对外开放80端口，防止暴露出特征，只有先使用默认账号连接到服务器在访问 [http://planet.cup](http://planet.cup)（位于manager/tcp.go这个文件） ，才可以访问控制面板。
默认账户信息:端口8366,密码12345678,加密方式AES-256-CFB，位于web-UI/webui/server.go文件，这个账户目前还不支持修改，后续版本将修复，如果有能力欢迎PR。
# Bugs

请保留 `go/src/github.com/Yee2` 路径下面的文件，目前还没有将模板登静态资源打包进程序里面。

# 配置文件说明

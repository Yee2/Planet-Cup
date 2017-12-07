package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"runtime"
	"fmt"
	"path/filepath"
	"flag"
	"io/ioutil"
	"encoding/json"
)

var config struct {
	Verbose    bool
	UDPTimeout time.Duration
}

func logf(f string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d: %s",filepath.Base(file),line,f)
	}
	if config.Verbose {
		log.Printf(f, v...)
	}
}

func main() {

	// 开启日志输出到文件

	config.Verbose = true
	//log.SetFlags( log.Ldate | log.Ltime | log.Lshortfile )
	//file,err := os.OpenFile("s.log",os.O_WRONLY | os.O_CREATE,0755)
	//if err != nil{
	//	log.Fatal("Failed to open log file!")
	//}
	//defer file.Close()
	//log.SetOutput(file)
	var configFile string
	flag.StringVar(&configFile,"c","config.json","Set the config file")

	f,err := os.Open(configFile)

	if err != nil{
		log.Fatalln("Failed to read configuration file")
	}

	str,err := ioutil.ReadAll(f)

	if err != nil{
		log.Fatalln(err)
	}
	configuration := struct {
		Verbose bool `json:"verbose"`
		Server string `json:"server"`
		Key string `json:"key"`
		Duration int64 `json:"duration"`
	}{}
	err = json.Unmarshal(str,&configuration)

	if err != nil {
		log.Fatalln(err)
	}
	if configuration.Server == ""{
		log.Fatalln("Server address can not be empty")
	}

	//开始运行
	if configuration.Duration > 0 {
		Duration = time.Second * time.Duration(configuration.Duration)
	}

	c,err := NewClient(configuration.Server ,configuration.Key )
	if err!=nil{
		log.Fatalln(err)
	}
	go c.run()


	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 停止监听
	c.shutdown()
}

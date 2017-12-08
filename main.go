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

	var (
		configFile string
		out string
	)
	flag.StringVar(&configFile,"c","config.json","Set the config file")
	flag.StringVar(&out,"o","","out file")
	flag.Parse()

	f,err := os.Open(configFile)

	if err != nil{
		log.Fatalf("Failed to read configuration file:%s",configFile)
	}
	defer f.Close()
	
	if out != ""{
		f_out,err := os.OpenFile(out,os.O_WRONLY | os.O_CREATE,0755)
		if err != nil{
			log.Fatalf("Failed to open file:%s",out)
		}
		defer f_out.Close()
		log.SetOutput(f_out)
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

	config.Verbose = configuration.Verbose

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

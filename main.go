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

	//开始运行
	c,err := NewClient("http://m2.local","b3tlqTsCs0Ipc0gW")
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

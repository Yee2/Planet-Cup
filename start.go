package main

//go:generate  go run generate/main.go

import (
	"github.com/Yee2/Planet-Cup/webui"
	"github.com/Yee2/Planet-Cup/ylog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {

	for i := range os.Args {
		if os.Args[i] == "--dev" {
			log.Println("developer mode")
			webui.BuiltIn = false
			go func() {
				log.Println(http.ListenAndServe(":6060", nil))
			}()
			break
		}
	}
	go ylog.Print()
	time.Sleep(time.Second)
	webui.Listen()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

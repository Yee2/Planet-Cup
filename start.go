package main

//go:generate  go run generate/main.go

import (
	"hyacinth/webui"
	"hyacinth/ylog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	for i := range os.Args {
		if os.Args[i] == "--dev" {
			webui.BuiltIn = false
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

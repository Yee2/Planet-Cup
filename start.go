package main

//go:generate  go run generate/main.go


import (
	"github.com/Yee2/Planet-Cup/webui"
	"github.com/Yee2/Planet-Cup/ylog"
	"os"
	"os/signal"
	"syscall"
)
func main()  {
	go ylog.Print()
	webui.Listen(8080)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

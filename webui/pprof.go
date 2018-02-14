package webui

import (
	_ "net/http/pprof"
	"log"
	"net/http"
)

func init()  {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
}

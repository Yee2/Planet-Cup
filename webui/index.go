package webui

import (
	"net/http"
	"github.com/Yee2/Planet-Cup/manager"
	"github.com/julienschmidt/httprouter"
)

var (
	shadowsocks_methods  = []string{"AES-128-GCM","AES-196-GCM","AES-256-GCM","AES-128-CTR",
		"AES-192-CTR","AES-256-CTR","AES-128-CFB","AES-192-CFB","AES-256-CFB","CHACHA20-IETF-POLY1305"}
)

func index(w http.ResponseWriter,r *http.Request, _ httprouter.Params){
	view_refresh(w,"index", struct {
		Methods interface{}
		List map[int]*manager.Shadowsocks
	}{Methods:shadowsocks_methods,List:tables.Rows})
}
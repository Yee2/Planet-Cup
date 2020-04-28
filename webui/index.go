package webui

import (
	"github.com/julienschmidt/httprouter"
	"hyacinth/manager"
	"net/http"
)

var (
	shadowsocks_methods = []string{"AES-128-GCM", "AES-196-GCM", "AES-256-GCM", "AES-128-CTR",
		"AES-192-CTR", "AES-256-CTR", "AES-128-CFB", "AES-192-CFB", "AES-256-CFB", "CHACHA20-IETF-POLY1305"}
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view(w, "index", struct {
		Methods interface{}
		List    map[int]*manager.Shadowsocks
	}{Methods: shadowsocks_methods, List: tables.Rows})
}

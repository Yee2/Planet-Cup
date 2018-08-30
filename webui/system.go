package webui

import (
	"code.cloudfoundry.org/bytefmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net"
	"net/http"
	"os"
	"runtime"
)

func system(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := make(map[string]interface{})
	info["逻辑CPU数量"] = runtime.NumCPU()
	info["Golang 版本"] = runtime.Version()
	info["Goroutine 数量"] = runtime.NumGoroutine()
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)
	info["程序占用内存"] = bytefmt.ByteSize(mem.Sys)
	if hostname, err := os.Hostname(); err == nil {
		info["主机名"] = hostname
	}
	if ifaces, err := net.Interfaces(); err == nil {
		var ip string
		// handle err
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				continue
			}
			// handle err
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					if v.IP.String() == "127.0.0.1" {
						continue
					}
					if v.IP.To4() != nil {
						ip += "<br>" + v.IP.String()
					}
				case *net.IPAddr:
					if v.IP.String() == "127.0.0.1" {
						continue
					}
					if v.IP.To4() != nil {
						ip += "<br>" + v.IP.String()
					}
				}
			}
		}
		info["IP 地址"] = template.HTML(ip)
	}
	view(w, "system", struct {
		Info interface{}
	}{info})
}

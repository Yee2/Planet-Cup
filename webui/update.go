package webui

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"strings"
	"strconv"
	"encoding/json"
	"time"
)

func update(w http.ResponseWriter,r *http.Request, ps httprouter.Params){

	port,err := strconv.Atoi(ps.ByName("port"))
	if err != nil{
		logger.Info("%s",err)
		http.NotFound(w,r)
		return
	}
	if port == 8366{
		res_error(w,http.StatusUnauthorized,"默认端口不能修改!")
		return
	}

	ss,ok := tables.Rows[port]
	if !ok{
		logger.Info("端口%d未找到",port)
		http.NotFound(w,r)
		return
	}

	encryption := strings.ToUpper(r.PostFormValue("encryption"))
	encryption_check := false
	for _,m := range shadowsocks_methods{
		if m == encryption{
			encryption_check = true
			break
		}
	}

	if !encryption_check{
		res_error(w,http.StatusBadRequest,"错误请求!")
		return
	}

	password :=  r.PostFormValue("password")
	if len(password) < 6{
		res_error(w,http.StatusBadRequest,"密码太短，最少6位")
		return
	}

	if len(password) > 32 {
		res_error(w,http.StatusBadRequest,"密码太长，最多32位")
		return
	}
	ss.Cipher = encryption
	ss.Password = password

	go func() {
		if err := ss.Stop(); err != nil{
			logger.Warning("无法正常关闭服务:%s",err)
			return
		}
		time.Sleep(time.Second * 5)//
		if err := ss.Start(); err != nil{
			logger.Warning("无法正常启动服务:%s",err)
			return
		}
	}()

	http.Error(w,http.StatusText(201),201)
	resp := json.NewEncoder(w)
	resp.Encode(struct {
		Port int `json:"port"`
		Password string `json:"password"`
		Method string `json:"method"`
	}{ss.Port,ss.Password,ss.Cipher})

}

package webui

import (
	"net/http"
	"strings"
	"strconv"
	"github.com/Yee2/Planet-Cup/manager"
	"github.com/julienschmidt/httprouter"
)

func add(w http.ResponseWriter,r *http.Request, _ httprouter.Params){
	method := strings.ToUpper(r.PostFormValue("method"))
	method_check := false
	for _,m := range shadowsocks_methods{
		if m == method{
			method_check = true
			break
		}
	}

	if !method_check{
		res_error(w,"错误请求!")
		return
	}

	port,err := strconv.Atoi(r.PostFormValue("port"))
	if err != nil{
		res_error(w,"请输入合法的端口范围")
		return
	}
	if port > 10000 || port < 1024{
		res_error(w,"端口只能设置在1024-10000之间")
		return
	}
	password :=  r.PostFormValue("password")
	if len(password) < 6{
		res_error(w,"密码太短，最少6位")
		return
	}

	if len(password) > 32 {
		res_error(w,"密码太长，最多32位")
		return
	}

	ss,err := manager.NewShadowsocks(port,password,method)
	if err != nil{
		res_error(w,"请确定输入信息无误。")
		return
	}


	if err := tables.Add(ss); err != nil {
		res_error(w,"添加失败，可能是端口被占用")
		return
	}
	res_message(w,"添加成功！")
}


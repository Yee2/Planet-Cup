package webui

import (
	"github.com/Yee2/Planet-Cup/manager"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
)

func add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	method := strings.ToUpper(r.PostFormValue("method"))
	methodCheck := false
	for _, m := range shadowsocks_methods {
		if m == method {
			methodCheck = true
			break
		}
	}

	if !methodCheck {
		res_error(w, 400, "错误请求!")
		return
	}

	port, err := strconv.Atoi(r.PostFormValue("port"))
	if err != nil {
		res_error(w, 400, "请输入合法的端口范围")
		return
	}
	if port > 10000 || port < 1024 {
		res_error(w, 400, "端口只能设置在1024-10000之间")
		return
	}
	password := r.PostFormValue("password")
	if len(password) < 6 {
		res_error(w, 400, "密码太短，最少6位")
		return
	}

	if len(password) > 32 {
		res_error(w, 400, "密码太长，最多32位")
		return
	}

	traffic, err := strconv.Atoi(r.PostFormValue("traffic"))
	if err != nil {
		res_error(w, 400, "请输入正确的流量限制。")
		return
	}

	ss, err := manager.NewShadowsocks(port, password, method)
	if err != nil {
		res_error(w, http.StatusBadRequest, "请确定输入信息无误。")
		return
	}
	ss.SetTrafficLimit(traffic)
	if err := tables.Add(ss); err != nil {
		res_error(w, http.StatusInternalServerError, "添加失败，可能是端口被占用")
		return
	}
	res_message(w, "添加成功！")
}

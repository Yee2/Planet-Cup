package webui

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func Password(w http.ResponseWriter,r *http.Request, _ httprouter.Params){
	if r.PostFormValue("old") != password{
		res_error(w,http.StatusBadRequest,"操作失败，密码错误")
		return
	}

	if r.PostFormValue("new") == ""{
		res_error(w,http.StatusBadRequest,"输入密码不能为空")
		return
	}
	if len(r.PostFormValue("new")) < 5 || len(r.PostFormValue("new")) > 32{
		res_error(w,http.StatusBadRequest,"设置的密码最小长度5，最大长度32")
		return
	}

	password = r.PostFormValue("new")
	logger.Info("修改管理密码为%s",password)
	res_message(w,"修改完成")
}

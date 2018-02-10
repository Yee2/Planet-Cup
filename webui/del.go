package webui

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func del(w http.ResponseWriter,r *http.Request, ps httprouter.Params){
	port,err := strconv.Atoi(ps.ByName("port"))
	if err != nil{
		res_error(w,"错误请求!")
		return
	}
	if err := tables.Del(port); err != nil{
		res_error(w,"删除失败，可能服务不错在!")
		return
	}
	res_message(w,"删除成功!")
}
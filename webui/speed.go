package webui

import (
	"golang.org/x/net/websocket"
	"time"
	"encoding/json"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func speed(w http.ResponseWriter,r *http.Request, ps httprouter.Params){
	port,err := strconv.Atoi(ps.ByName("port"))
	if err != nil{
		http.NotFound(w,r)
		return
	}
	ss,ok := tables.Rows[port]
	if !ok{
		http.NotFound(w,r)
		return
	}
	websocket.Handler(func(ws *websocket.Conn) {
		type data struct {
			X int `json:"up"`
			Y int `json:"down"`
			Time int64 `json:"time"`
		}
		history := make([]data,0)
		resp := json.NewEncoder(ws)
		for _,v := range ss.Speed_history(30){
			history = append(history,data{v.Up,v.Down,v.Time})
		}
		resp.Encode(history)
		for {
			<- time.Tick(time.Second)
			up,down := ss.Speed()
			resp.Encode([]data{{up,down,time.Now().Unix()}})
		}
	}).ServeHTTP(w,r)
}
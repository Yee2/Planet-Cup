package webui

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/Yee2/Planet-Cup/ylog"
	"sort"
)

func logs(w http.ResponseWriter,r *http.Request, _ httprouter.Params){
	pool := ylog.GetAll()
	Logs := make(log_list,0)
	for name := range pool{
		for _,log := range pool[name].All(){
			Logs = append(Logs,struct{
				Content string
				Date int64
				Level string
				File string
				Class string
				Line int
			}{log.Content,log.Time.Unix(),log.Level.String(),log.File,name,log.Line})
		}
	}
	sort.Sort(Logs)
	view_refresh(w,"logger", struct {
		Logs interface{}
	}{Logs})
}

type log_list []struct{
	Content string
	Date int64
	Level string
	File string
	Class string
	Line int
}

func (self log_list)Len()int  {
	return len(self)
}
func (self log_list)Less(i , j int)bool  {
	return self[i].Date < self[j].Date
}
func (self log_list)Swap(i , j int) {
	self[i],self[j] = self[j],self[i]
}
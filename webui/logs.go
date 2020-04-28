package webui

import (
	"github.com/julienschmidt/httprouter"
	"hyacinth/ylog"
	"net/http"
	"sort"
)

func logs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pool := ylog.GetAll()
	Logs := make(logList, 0)
	for name := range pool {
		for _, log := range pool[name].All() {
			Logs = append(Logs, struct {
				Content string
				Date    int64
				Level   string
				File    string
				Class   string
				Line    int
			}{log.Content, log.Time.Unix(), log.Level.String(), log.File, name, log.Line})
		}
	}
	sort.Sort(Logs)
	view(w, "logger", struct {
		Logs interface{}
	}{Logs})
}

type logList []struct {
	Content string
	Date    int64
	Level   string
	File    string
	Class   string
	Line    int
}

func (list logList) Len() int {
	return len(list)
}
func (list logList) Less(i, j int) bool {
	return list[i].Date < list[j].Date
}
func (list logList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

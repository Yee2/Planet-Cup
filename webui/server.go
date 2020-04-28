package webui

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry/bytefmt"
	"github.com/fsnotify/fsnotify"
	"github.com/julienschmidt/httprouter"
	"hyacinth/manager"
	"hyacinth/ylog"
)

var (
	view_entry *template.Template
	views      map[string]*template.Template
	err        error
)

var mux = &http.ServeMux{}
var logger = ylog.NewLogger("web-UI")
var tables = manager.NewTable()
var base = template.New("base").Funcs(template.FuncMap{"ByteSize": ByteSize, "Date": date})
var BuiltIn = true
var initial = false

func viewInitial() {
	logger.Info("Initialize template resources")
	views = make(map[string]*template.Template)
	if BuiltIn {
		logger.Info("Use built-in template resources")
		// 使用内置的资源
		assets, err := NewAssets(TemplateAssets[:])
		letItDie(err)
		if resource, ok := assets["/entry.html"]; ok {
			data, err := ioutil.ReadAll(resource)
			letItDie(err)
			view_entry = template.Must(base.Parse(string(data)))
		} else {
			panic(errors.New("unable to read the entry file"))
		}
		// 加载全部组件
		for name := range assets {
			if strings.HasPrefix(name, "/components/") && strings.HasSuffix(name, ".html") {
				data, err := ioutil.ReadAll(assets[name])
				letItDie(err)
				view_entry = template.Must(base.Parse(string(data)))
			}
		}

		for name := range assets {
			if strings.HasPrefix(name, "/content/") && strings.HasSuffix(name, ".html") {
				data, err := ioutil.ReadAll(assets[name])
				letItDie(err)
				views[filepath.Base(name)] = template.Must(template.Must(view_entry.Clone()).Parse(string(data)))
			}
		}

	} else {
		// 使用本地路径资源
		refresh()
		daemon()

	}
}
func daemon() {
	watch, err := fsnotify.NewWatcher()
	letItDie(err)
	letItDie(watch.Add("assets/template"))

	dida := 0
	t := time.NewTicker(time.Second)
	changed := false
	sign := make(chan int)
	stop := make(chan int)
	go func() {
		for {
			select {
			case <-t.C:
				if !changed {
					continue
				}
				dida++
				if dida > 3 {
					refresh()
					changed = false
				}
			case <-sign:
				changed = true
				dida = 0
			case <-stop:
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-watch.Events:
				{
					sign <- 1
				}
			case err := <-watch.Errors:
				{
					if err != nil {
						logger.Danger("%s", err)
						stop <- 1
						return
					}
				}
			}
		}
		watch.Close()
	}()
}
func refresh() {
	entry := template.Must(base.ParseFiles("assets/template/entry.html"))
	entry, err = view_entry.ParseGlob("assets/template/components/*.html")
	if err != nil {
		logger.Warning("%s", err)
		return
	}
	view_entry = entry
	files, err := filepath.Glob("assets/template/content/*.html")
	letItDie(err)
	for _, f := range files {
		t, err := template.Must(view_entry.Clone()).ParseFiles(f)
		if err != nil {
			logger.Warning("%s", err)
			continue
		}
		views[filepath.Base(f)] = t
	}
}
func Listen() {
	err := tables.Load("data.json")
	if err != nil {
		logger.Danger("%s", err)
	}
	tables.Boot()
	router := httprouter.New()
	// 免登陆部分
	router.GET("/login.html", login)
	router.GET("/logout.html", logout)
	router.POST("/login.html", loginVerify)

	if BuiltIn {
		FileHandle, err := NewAssets(AssetsData[:])
		if err != nil {
			panic(err)
		}
		logger.Info("Use file built-in resources")
		router.ServeFiles("/public/*filepath", FileHandle)
	} else {
		logger.Info("Use external static resources")
		router.ServeFiles("/public/*filepath", http.Dir("assets/public"))
	}

	// 登录后可看部分
	router.GET("/", auth(index))
	router.GET("/index.html", auth(index))
	router.GET("/system.html", auth(system))
	router.GET("/logger.html", auth(logs))
	router.POST("/password.html", auth(Password))

	router.GET("/shadowsocks/:port/speed", auth(speed))
	router.POST("/shadowsocks", auth(add))
	router.PUT("/shadowsocks/:port", auth(update))
	router.DELETE("/shadowsocks/:port", auth(del))

	go func() {
		for {
			time.Sleep(time.Minute * 10)
			err := tables.Save("data.json")
			if err != nil {
				logger.Danger("saving data failed:%s", err)
			}
		}
	}()
	go func() {
		err := http.ListenAndServe("0.0.0.0:34567", router)
		if err != nil {
			logger.Info("initializing the web server failed:%s", err)
		}
	}()
}

func letItDie(err error) {
	if err != nil {
		panic(err)
	}
}
func view(w io.Writer, name string, data interface{}) {

	if initial == false {
		initial = true
		viewInitial()
	}
	if tpl, ok := views[name+".html"]; ok {
		tpl.ExecuteTemplate(w, "entry", data)
		return
	}

	panic(fmt.Errorf("view does not exist：%s", name+".html"))
}

func res_error(w http.ResponseWriter, code int, text string) {
	w.Header().Set("Content-Type", "Content-Type: application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(code)
	response := json.NewEncoder(w)
	response.Encode(struct {
		Error string `json:"error"`
		Code  int    `json:"code"`
	}{text, -1})
}
func res_message(w io.Writer, text string) {
	response := json.NewEncoder(w)
	response.Encode(struct {
		Error string `json:"message"`
		Code  int    `json:"code"`
	}{text, 0})
}
func ByteSize(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	if size, ok := args[0].(int); ok {
		return bytefmt.ByteSize(uint64(size))
	}

	if str, ok := args[0].(string); ok {
		if size, err := strconv.Atoi(str); err == nil {
			return bytefmt.ByteSize(uint64(size))
		}
	}
	return ""
}
func date(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	if timestamp, ok := args[0].(int64); ok {
		t := time.Unix(timestamp, 0)
		return t.Format(time.ANSIC)
	}
	return ""
}

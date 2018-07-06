package webui

import (
	"net/http"
	"html/template"
	"path/filepath"
	"io"
	"errors"
	"strconv"
	"github.com/Yee2/Planet-Cup/ylog"
	"encoding/json"
	"github.com/Yee2/Planet-Cup/manager"
	"code.cloudfoundry.org/bytefmt"
	"github.com/julienschmidt/httprouter"
	"time"
	"io/ioutil"
	"strings"
	"fmt"
)

var (
	view_entry *template.Template
	views      map[string]*template.Template
	err        error
)

// 静态资源部分
var AssetsData []byte
var TemplateAssets []byte

var mux = &http.ServeMux{}
var logger = ylog.NewLogger("web-UI")
var tables = manager.NewTable()
var base = template.New("base").Funcs(template.FuncMap{"ByteSize": ByteSize, "Date": date})
var initial = false
func viewInitial() {
	logger.Info("初始化模板资源...")
	views = make(map[string]*template.Template)
	if len(TemplateAssets) > 0 {
		logger.Info("使用内置模板资源...")
		// 使用内置的资源
		assets, err := NewAssets(TemplateAssets)
		letItDie(err)
		if resource, ok := assets["/entry.html"]; ok {
			data,err := ioutil.ReadAll(resource)
			letItDie(err)
			view_entry = template.Must(base.Parse(string(data)))
		}else{
			panic(errors.New("无法读取入口文件"))
		}
		// 加载全部组件
		for name := range assets{
			if strings.HasPrefix(name,"/components/") && strings.HasSuffix(name,".html"){
				data,err := ioutil.ReadAll(assets[name])
				letItDie(err)
				view_entry = template.Must(base.Parse(string(data)))
			}
		}

		for name := range assets{
			if strings.HasPrefix(name,"/content/") && strings.HasSuffix(name,".html"){
				data,err := ioutil.ReadAll(assets[name])
				letItDie(err)
				views[filepath.Base(name)] = template.Must(template.Must(view_entry.Clone()).Parse(string(data)))
			}
		}

	} else {
		// 使用本地路径资源
		view_entry = template.Must(base.ParseFiles("assets/template/entry.html"))
		view_entry, err = view_entry.ParseGlob("assets/template/components/*.html")
		letItDie(err)
		files, err := filepath.Glob("assets/template/content/*.html")
		letItDie(err)
		for _, f := range files {
			views[filepath.Base(f)] = template.Must(template.Must(view_entry.Clone()).ParseFiles(f))
		}
	}
}

func Listen() {
	ss, _ := manager.NewShadowsocks(8366, "12345678", "AES-256-CFB")
	logger.Info("启动ShadowSock服务")
	tables.Push(ss)
	tables.Boot()
	router := httprouter.New()
	// 免登陆部分
	router.GET("/login.html", login)
	router.GET("/logout.html", logout)
	router.POST("/login.html", login_chick)

	if len(AssetsData) > 0 {

		FileHandle, err := NewAssets(AssetsData)
		if err != nil {
			panic(err)
		}
		logger.Info("使用文件内置资源")
		router.ServeFiles("/public/*filepath", FileHandle)
	} else {
		logger.Info("使用外置静态资源")
		router.ServeFiles("/public/*filepath", http.Dir("assets/public"))
	}

	// 登录后可看部分
	router.GET("/", auth(index))
	router.GET("/index.html", auth(index))
	router.GET("/system.html", auth(system))
	router.GET("/logger.html", auth(logs))
	router.GET("/shadowsocks/:port/speed", auth(speed))
	router.POST("/shadowsocks", auth(add))
	router.PUT("/shadowsocks/:port", auth(update))
	router.DELETE("/shadowsocks/:port", auth(del))

	err := http.ListenAndServe("0.0.0.0:34567", router)
	if err != nil {
		logger.Info("初始化Web服务器失败:%s", err)
	}
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

	panic(fmt.Errorf("视图不存在：%s",name+".html"))
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

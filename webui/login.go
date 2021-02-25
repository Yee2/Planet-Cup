package webui

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var session = struct {
	rows map[string]time.Time //表示过期时间
	sync.Mutex
}{make(map[string]time.Time), sync.Mutex{}}
var (
	user     = "admin"
	password = "admin"
)

func init() {
	go func() {
		for {
			// 清理的过期的会话
			time.Sleep(time.Hour)
			for i, v := range session.rows {
				if v.Before(time.Now()) {
					session.Lock()
					delete(session.rows, i)
					session.Unlock()
				}
			}
		}
	}()
}
func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if isLogin(w, r) {
		http.Redirect(w, r, "/index.html", 301)
		return
	}
	if BuiltIn {
		t := template.Must(template.ParseFS(templates, "assets/template/components/head.html", "assets/template/login.html"))
		t.ExecuteTemplate(w, "login", nil)
	} else {
		t := template.Must(template.ParseFiles("assets/template/components/head.html", "assets/template/login.html"))
		t.ExecuteTemplate(w, "login", nil)
	}
}

func loginFail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if isLogin(w, r) {
		http.Redirect(w, r, "/index.html", 301)
		return
	}
	if BuiltIn {
		t := template.Must(template.ParseFS(templates, "assets/template/components/head.html", "assets/template/login.html"))
		t.ExecuteTemplate(w, "login", struct {
			Message string
		}{"用户名或密码错误，请重试。"})
	} else {
		t := template.Must(template.ParseFiles("assets/template/components/head.html", "assets/template/login.html"))
		t.ExecuteTemplate(w, "login", struct {
			Message string
		}{"用户名或密码错误，请重试。"})
	}
}

func loginVerify(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.PostFormValue("user") == user && r.PostFormValue("password") == password {
		key := RandStr(8)
		session.Lock()
		defer session.Unlock()
		session.rows[key] = time.Now().Add(time.Minute * 10) //有效期十分钟
		cookie := http.Cookie{Name: "session_key", Value: key}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/index.html", 301)
		return
	}
	loginFail(w, r, ps)
}
func logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie, err := r.Cookie("session_key")
	if err == nil {
		if _, ok := session.rows[cookie.Value]; ok {
			session.Lock()
			delete(session.rows, cookie.Value)
			session.Unlock()
		}
	}
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login.html", 301)
}
func isLogin(_ http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("session_key")
	if err != nil {
		return false
	}
	if s, ok := session.rows[cookie.Value]; ok {
		if s.After(time.Now()) {
			session.Lock()
			defer session.Unlock()
			s = time.Now().Add(time.Minute * 10)
			return true
		}
	}
	return false
}

func auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if isLogin(w, r) {
			h(w, r, ps)
		} else {
			w.Header().Set("Cache-Control", "no-cache")
			http.Redirect(w, r, "/login.html", 301)
		}
	}
}

//生成随机字符串
func RandStr(strlen int) string {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, strlen)
	var num int
	for i := 0; i < strlen; i++ {
		num = rand.Intn(57) + 65
		for {
			if num > 90 && num < 97 {
				num = rand.Intn(57) + 65
			} else {
				break
			}
		}
		data[i] = byte(num)
	}
	return string(data)
}

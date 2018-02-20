// 这是一个对golang官方log包的包装，将日志写入部分和读取部分分开管理
// 注意：这个包只做日志记录，不会停止程序运行，更不会丢出panic(初始化除外)
// 这个包不会返回任何错误
package ylog

import (
	"time"
	"runtime"
	"fmt"
	"errors"
	"sync"
)

type Log struct {
	Time time.Time
	File string
	Line int
	Content string
	Level Level
}

type Level int

func (self Level) String() string {
	switch self {
	case DANGER:
		return "danger"
	case SUCCESS:
		return "success"
	case WARNING:
		return "warning"
	case INFO:
		return "info"
	default:
		return "notice"
	}
}
const (
	_ Level  = iota
	SUCCESS
	INFO
	WARNING
	DANGER
)

var (
	logs = struct{
		pool map[string]*Logger
		sync.Mutex
	}{make(map[string]*Logger),sync.Mutex{}}
)
func GetAll() map[string]*Logger {
	return logs.pool
}

type Logger struct {
	key string
	offset int
	list [1024]*Log
	sync.Mutex
}

func NewLogger(name string) *Logger{
	if name == ""{
		panic(errors.New("name is empty"))
	}
	if _,ok := logs.pool[name]; ok{
		panic(errors.New("this name already exists"))
	}
	logger := &Logger{key:name,offset:0,list:[1024]*Log{}}
	logs.Lock()
	defer logs.Unlock()
	logs.pool[name] = logger
	return logger
}
func (self *Logger)next(){
	self.offset ++
	if self.offset >= 1024{
		self.offset -= 1024
	}
}
func (self *Logger)All()[]*Log{
	r := make([]*Log,0,1024)
	for i := 1023; i >= 0 ;i --{
		key := i + self.offset
		if key >= 1024{
			key -= 1024
		}
		if self.list[key] == nil{
			break
		}
		r = append(r,self.list[key])
	}
	return  r
}

func (self *Logger)insert(level Level,Content string, agrs []interface{})  {
	_, file, line, _ := runtime.Caller(2)
	Content = fmt.Sprintf(Content, agrs ...)
	log := &Log{
		Time:time.Now(),
		File:file,
		Line:line,
		Content:Content,
		Level:level,
	}
	self.Lock()
	defer self.Unlock()
	self.list[self.offset] = log
	go broadcast(log)//广播最新日志
	self.next()
}
func (self *Logger)Info(Content string, agrs ... interface{})  {
	self.insert(INFO,Content,agrs)
}
func (self *Logger)Success(Content string, agrs ... interface{})  {
	self.insert(SUCCESS,Content,agrs)
}
func (self *Logger)Warning(Content string, agrs ... interface{})  {
	self.insert(WARNING,Content,agrs)
}
func (self *Logger)Danger(Content string, agrs ... interface{})  {
	self.insert(DANGER,Content,agrs)
}
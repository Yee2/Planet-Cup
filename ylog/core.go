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

const (
	_ Level  = iota
	SUCCESS
	INFO
	WARNING
	DANGER
)

var (
	logs = struct{
		pool map[string][1024]*Log
		*sync.Mutex
	}{make(map[string][1024]*Log),new(sync.Mutex)}

	listener chan *Log
)

type Logger struct {
	key string
	offset int
}

func NewLogger(name string) *Logger{
	if name == ""{
		panic(errors.New("name is empty"))
	}
	if _,ok := logs.pool[name]; ok{
		panic(errors.New("this name already exists"))
	}
	logger := &Logger{name,0}
	logs.Lock()
	defer logs.Unlock()
	logs.pool[name] = [1024]*Log{}
	return logger
}
func (self *Logger)next(){
	self.offset ++
	if self.offset >= 1024{
		self.offset -= 1024
	}
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
	logs.Lock()
	defer logs.Unlock()
	if logger,ok := logs.pool[self.key]; ok{
		logger[self.offset] = log
	}else{
		logs.pool[self.key] = [1024]*Log{log}
	}
	if listener != nil{
		listener <- log
	}
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
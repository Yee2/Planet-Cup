// 这是一个对golang官方log包的包装，将日志写入部分和读取部分分开管理
// 注意：这个包只做日志记录，不会停止程序运行，更不会丢出panic(初始化除外)
// 这个包不会返回任何错误
package ylog

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Log struct {
	Time    time.Time
	File    string
	Line    int
	Content string
	Level   Level
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
	_ Level = iota
	SUCCESS
	INFO
	WARNING
	DANGER
)

var (
	logs = struct {
		pool map[string]*Logger
		sync.Mutex
	}{make(map[string]*Logger), sync.Mutex{}}
)

func GetAll() map[string]*Logger {
	return logs.pool
}

type Logger struct {
	key    string
	offset int
	list   [1024]*Log
	sync.Mutex
}

func NewLogger(name string) *Logger {
	if name == "" {
		panic(errors.New("name is empty"))
	}
	if _, ok := logs.pool[name]; ok {
		panic(errors.New("this name already exists"))
	}
	logger := &Logger{key: name, offset: 0, list: [1024]*Log{}}
	logs.Lock()
	defer logs.Unlock()
	logs.pool[name] = logger
	return logger
}
func (logger *Logger) next() {
	logger.offset++
	if logger.offset >= 1024 {
		logger.offset = 0
	}
}
func (logger *Logger) All() []*Log {
	r := make([]*Log, 0, 1024)
	for i := 1023; i >= 0; i-- {
		key := i + logger.offset
		if key >= 1024 {
			key -= 1024
		}
		if logger.list[key] == nil {
			break
		}
		r = append(r, logger.list[key])
	}
	return r
}

func (logger *Logger) insert(level Level, Content string, args []interface{}) {
	_, file, line, _ := runtime.Caller(2)
	Content = fmt.Sprintf(Content, args...)
	log := &Log{
		Time:    time.Now(),
		File:    file,
		Line:    line,
		Content: Content,
		Level:   level,
	}
	logger.Lock()
	defer logger.Unlock()
	logger.list[logger.offset] = log
	//广播最新日志
	for _, v := range receivers.list {
		select {
		case v.ch <- log:
		default:
		}
	}
	logger.next()
}
func (logger *Logger) Info(Content string, args ...interface{}) {
	logger.insert(INFO, Content, args)
}
func (logger *Logger) Success(Content string, args ...interface{}) {
	logger.insert(SUCCESS, Content, args)
}
func (logger *Logger) Warning(Content string, args ...interface{}) {
	logger.insert(WARNING, Content, args)
}
func (logger *Logger) Danger(Content string, args ...interface{}) {
	logger.insert(DANGER, Content, args)
}

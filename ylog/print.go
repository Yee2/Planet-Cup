package ylog

import (
	"fmt"
	"path/filepath"
	"strconv"
)

func Print() {
	listener := NewReceiver()
	defer listener.Close()
	for {
		log, err := listener.Receive()
		if err != nil {
			break
		}
		var level string
		t := NewTxt(fmt.Sprintf("[%s#%s] %s#%d:",
			level, log.Time.Format("01-02 15:04:05"),
			filepath.Base(log.File), log.Line))
		switch log.Level {
		case INFO:
			level = "info"
			t.Color(Blue)
		case WARNING:
			level = "warning"
			t.Color(YELLOW)
		case SUCCESS:
			level = "success"
			t.Color(GREEN)
		case DANGER:
			level = "danger"
			t.Color(RED)
		}
		fmt.Printf("%s%s\n", t, log.Content)
	}
}

type Color int

const (
	BLACK Color = iota
	RED
	GREEN
	YELLOW
	Blue
	MEGENTA
	CYAN
	WHITE
)

type Mode string

const (
	RESET      Mode = "0" //关闭所有属性
	BRIGHT          = "1" // 高亮度
	DIM             = "2"
	UNDER_LINE      = "4" // 下划线
	BLINK           = "5" // 闪烁
	REVERSE         = "7" // 反显示
	HIDDEN          = "8" //消隐
)

type txt struct {
	text,      //文字
	backColor, //背景色
	mode,      //控制符
	color string //文字颜色
}

func NewTxt(Content string) *txt {
	return &txt{text: Content, mode: "01"}
}

func (this *txt) Color(c Color) {
	this.color = "3" + strconv.Itoa(int(c))
}
func (this *txt) BackColor(c Color) {
	this.backColor = "4" + strconv.Itoa(int(c))
}
func (this *txt) Mode(m Mode) {
	this.mode = string(m)
}
package ylog

import (
	"fmt"
	"path/filepath"
	"strconv"
	"bytes"
)
func Print()  {
	listener := NewReceiver()
	defer listener.Close()
	for {
		log,err := listener.Receive()
		if err !=nil{
			break
		}
		var level string
		switch log.Level {
		case INFO:
			level = "info"
		case WARNING:
			level = "warning"
		case SUCCESS:
			level = "success"
		case DANGER:
			level = "danger"
		}
		fmt.Printf("[%s#%s] %s#%d: %s\n",level,log.Time,filepath.Base(log.File),log.Line,log.Content)
	}
}

const STR_BEGIN = "\033["
const STR_END = "\033[0m"

type Color int
const (
	BLACK Color= iota
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
	RESET   Mode    = "0" //关闭所有属性
	BRIGHT     = "1" // 高亮度
	DIM        = "2"
	UNDER_LINE = "4" // 下划线
	BLINK      = "5" // 闪烁
	REVERSE    = "7" // 反显示
	HIDDEN     = "8" //消隐
)

type txt struct {
	text,//文字
	backColor,//背景色
	mode,//控制符
	color string//文字颜色
}
func NewTxt(Content string) *txt{
	return &txt{text:Content,mode:"0"}
}

func (this *txt)Color(c Color){
	this.color = "3" + strconv.Itoa(int(c))
}
func (this *txt)BackColor(c Color){
	this.backColor = "4" + strconv.Itoa(int(c))
}
func (this *txt)Mode(m Mode){
	this.mode = string(m)
}
func (this *txt)ToString() string{
	buf := &bytes.Buffer{}
	buf.WriteString(STR_BEGIN)
	buf.WriteString(this.mode + ";")
	if this.color != "" {
		buf.WriteString(this.color + ";")
	}
	if this.backColor != "" {
		buf.WriteString(this.backColor + ";")
	}

	buf.WriteString("m")
	buf.WriteString(this.text)
	buf.WriteString(STR_BEGIN)
	return buf.String()
}

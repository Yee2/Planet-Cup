package ylog

import "bytes"

func (this *txt) String() string {
	buf := &bytes.Buffer{}
	buf.WriteByte(0x1B)
	buf.WriteByte('[')
	buf.WriteString(this.mode)
	if this.color != "" {
		buf.WriteString(";"+this.color)
	}
	if this.backColor != "" {
		buf.WriteString(";"+this.backColor)
	}

	buf.WriteString(`m`)
	buf.WriteString(this.text)
	buf.WriteByte(0x1B)
	buf.WriteString("[00m")
	return buf.String()
}


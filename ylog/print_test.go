package ylog

import "testing"

func TestPrint(t *testing.T) {
	text := NewTxt("Test...")
	text.Color(RED)
	t.Logf("\n%s\n", text)
}

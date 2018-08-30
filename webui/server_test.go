package webui

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"testing"
)

func TestTemplate(t *testing.T) {
	assets, err := NewAssets(TemplateAssets[:])
	if err != nil {
		t.Fatal(err)
	}
	data_head, err := ioutil.ReadAll(assets["/components/head.html"])
	if err != nil {
		t.Fatal(err)
	}
	data_login, err := ioutil.ReadAll(assets["/login.html"])
	if err != nil {
		t.Fatal(err)
	}
	buffer := bytes.NewBuffer([]byte{})
	view_login = template.Must(template.Must(template.New("").Parse(string(data_head))).Parse(string(data_login)))
	view_login.ExecuteTemplate(buffer, "login", nil)
	if buffer.Len() < 1 {
		t.Fatal()
	}
	t.Logf("============\n%s\n===============\n", buffer.Bytes())
}

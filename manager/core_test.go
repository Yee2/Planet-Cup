package manager

import (
	"io/ioutil"
	"net"
	"path/filepath"
	"testing"
)

func TestTable_Save(t *testing.T) {
	listen, err := net.Listen("tcp", "0.0.0.0:8300")
	if err != nil {
		t.Skipf("端口被占用，跳过本次测试")
	}

	listen.Close()

	table := NewTable()
	ss, err := NewShadowsocks(8300, "12345678", "AES-256-GCM")
	if err != nil {
		t.Fatal(err)
	}

	err = table.Add(ss)
	if err != nil {
		t.Fatal(err)
	}

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	table.Rows[8300].Stop()
	table.Save(filepath.Join(dir, "tmp.json"))
	t.Run("load", func(t *testing.T) {
		testTable_Save(t, filepath.Join(dir, "tmp.json"))
	})
}

func testTable_Save(t *testing.T, file string) {
	table := NewTable()
	err := table.Load(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Rows) < 1 {
		t.Fatalf("无法加载shadowsocks data")
	}
}

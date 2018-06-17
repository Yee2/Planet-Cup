package webui

import (
	"net/http"
	"bytes"
	"os"
	"errors"
	"archive/tar"
	"strings"
	"io"
)
var FileHandle Assets
func init()  {
	FileHandle = make(Assets)
	r := bytes.NewReader(AssetsData)
	tr := tar.NewReader(r)


	for {
		hdr,err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		data := make([]byte,hdr.Size)
		_,err = tr.Read(data)
		if err != io.EOF && err != nil{
			panic(err)
		}
		path := hdr.Name
		if !strings.HasPrefix(path,"/"){
			path = "/" + hdr.Name
		}
		FileHandle[path] = &File{Reader:bytes.NewReader(data),FileInfo:hdr.FileInfo()}


	}
}

type Assets map[string]*File

func (asset Assets)Open(name string) (http.File,error) {
	if f,ok := asset[name]; ok{
		return f,nil
	}
	logger.Warning("文件未找到:%s",name)
	return nil,errors.New("文件未找到")
}
type File struct {
	*bytes.Reader
	os.FileInfo
	err error
}

func (f *File)Readdir(count int) ([]os.FileInfo, error)  {
	return nil,errors.New("暂时不支持该操作!")
}


func (f *File)Close() error {
	//f = nil
	return nil
}

func (f *File)Stat() (os.FileInfo, error) {
	return f.FileInfo,f.err
}


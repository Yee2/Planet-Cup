package main

import (
	"os"
	"archive/tar"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"strings"
	"io"
	"bufio"
)

type Resources struct {
	w io.WriteCloser
	i uint16
}

func (r *Resources) Write(p []byte) (n int, err error) {
	w := bufio.NewWriter(r.w)
	defer w.Flush()
	for i := range p {
		_, err = w.WriteString(fmt.Sprintf("0x%02x,", p[i]))
		n ++
		if err != nil {
			break
		}else{
			r.i ++
			if r.i%20 == 0{
				w.WriteString("\n")
			}
		}
	}
	return
}
func (r *Resources) Close() error {
	r.w.Write([]byte{'}'})
	return r.w.Close()
}
func R(file, name string) (r *Resources, err error) {
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	fd.Truncate(0)
	fd.WriteString(fmt.Sprintf(`package webui
var %s = []byte{
`, name))
	return &Resources{fd,0}, nil
}

func main()  {
	tarfs,err := R("webui/AssetsData.go","AssetsData")
	if err != nil{
		fmt.Println(err)
		return
	}
	writer := tar.NewWriter(tarfs)
	defer writer.Flush()
	defer writer.Close()
	defer tarfs.Close()
	err = file2tar(writer,"assets/public/","assets/public/")
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println("操作完成!\n")
	}
}

func file2tar(w *tar.Writer,basedir,name string) (e error ){
	defer func() {
		if err := recover(); err != nil{
			e = err.(error)
		}
	}()
	files,err := ioutil.ReadDir(name)
	if err != nil{
		return err
	}
	for _,file := range files{
		if file.IsDir(){
			file2tar(w,basedir,filepath.Join(name,file.Name()))
		}else{
			info,err := tar.FileInfoHeader(file,"")
			letItDie(err)

			info.Name = strings.TrimPrefix(filepath.Join(name,file.Name()) ,basedir)
			fmt.Println("添加文件:",info.Name)
			err = w.WriteHeader(info)
			letItDie(err)
			f,err := os.Open(filepath.Join(name,file.Name()))
			letItDie(err)
			data,err := ioutil.ReadAll(f)
			letItDie(err)
			_,err = w.Write(data)
			letItDie(err)
		}
	}
	return nil
}

func letItDie(err error)  {
	if err != nil{
		panic(err)
	}
}
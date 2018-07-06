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
	i uint16 //写入字节数
}

func (r *Resources) Write(p []byte) (n int, err error) {
	w := bufio.NewWriter(r.w)
	defer w.Flush()
	for i := range p {
		_, err = w.WriteString(fmt.Sprintf("0x%02x,", p[i]))
		n ++
		if err != nil {
			break
		} else {
			r.i ++
			if r.i%20 == 0 {
				w.WriteString("\n")
			}
		}
	}
	return
}
func (r *Resources) Close() error {
	r.w.Write([]byte("}\r\n}"))
	return r.w.Close()
}
func R(file, name string) (r *Resources, err error) {
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	fd.Truncate(0)
	fd.WriteString(fmt.Sprintf(`package webui
func init(){
	%s = []byte{
`, name))
	return &Resources{fd, 0}, nil
}

func main() {
	f1()
	f2()
}
func f1()  {
	tarfs, err := R("webui/AssetsData.go", "AssetsData")
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := tar.NewWriter(tarfs)
	defer writer.Flush()
	defer writer.Close()
	defer tarfs.Close()
	err = file2tar(writer, "assets/public/", "assets/public/")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("操作完成!")
	}
}
func f2()  {
	tarfs, err := R("webui/TemplateAssets.go", "TemplateAssets")
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := tar.NewWriter(tarfs)
	defer writer.Flush()
	defer writer.Close()
	defer tarfs.Close()
	err = file2tar(writer, "assets/template/", "assets/template/")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("操作完成!")
	}
}
func file2tar(w *tar.Writer, basedir, name string) (e error) {
	name = Windows2Linux(name)
	defer func() {
		if err := recover(); err != nil {
			e = err.(error)
		}
	}()
	files, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}
	fmt.Printf("压缩%s",name)
	for _, file := range files {
		if file.IsDir() {
			fmt.Println()
			file2tar(w, basedir, filepath.Join(name, file.Name()))
		} else {
			info, err := tar.FileInfoHeader(file, "")
			letItDie(err)
			info.Name = strings.TrimPrefix(Windows2Linux(filepath.Join(name, file.Name())), basedir)
			fmt.Print(".")
			err = w.WriteHeader(info)
			letItDie(err)
			f, err := os.Open(filepath.Join(name, file.Name()))
			letItDie(err)
			data, err := ioutil.ReadAll(f)
			letItDie(err)
			_, err = w.Write(data)
			letItDie(err)
		}
	}
	fmt.Println()
	return nil
}

func letItDie(err error) {
	if err != nil {
		panic(err)
	}
}

func Windows2Linux(s string)(string){
	return strings.Replace(s,"\\","/",-1)
}
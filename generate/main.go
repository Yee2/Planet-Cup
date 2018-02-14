package main

import (
	"fmt"
	"os"
)

type Directory struct {
	Name string
	Files []*File
	Dirs []*Directory
}

type File struct {
	Name string
	Data []byte
}
func main()  {
	fmt.Printf("%+v\n%+v",os.Args,os.Environ())
}

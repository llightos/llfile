package main

import (
	"fmt"
	"os"
)

//返回下载的断点
func Length(fileName string) int64 {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	fmt.Println(stat.Size())
	return stat.Size()
}

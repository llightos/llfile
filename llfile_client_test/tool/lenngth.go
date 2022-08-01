package tool

import (
	"fmt"
	"os"
)

//返回下载的断点
func Length(fileName string) int64 {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0666)
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

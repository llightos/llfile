package tool

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"
)

func OpenFile(way string) (*os.File, os.FileInfo) {
	file, err := os.Open(way)
	if err != nil {
		panic(err)
	}
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	return file, stat
}

func FileMd5(file *os.File) string {
	md5 := md5.New()
	_, err := io.Copy(md5, file)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", md5.Sum(nil))
}

func DivideName(name string) (headName string, expandName string) {
	split := strings.Split(name, ".")
	return split[0], split[1]
}

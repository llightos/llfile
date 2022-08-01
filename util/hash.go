package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

// FileHash 判断file的md5值与str等不等
func FileHash(str string) bool {
	file, err2 := os.Open("./file/" + str + ".llfile")
	if err2 != nil {
		log.Println("FileHash(str string) err")
		return false
	}
	hash := md5.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		fmt.Println("FileHash(file *os.File, str string), 错误")
		return false
	}
	sum := hash.Sum(nil)
	if fmt.Sprintf("%x", sum) == str {
		return true
	}
	return false
}

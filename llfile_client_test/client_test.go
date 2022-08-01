package main

import (
	"fmt"
	"llfile_client/tool"
	"os"
	"testing"
)

func TestFileMd5(t *testing.T) {
	open, err := os.Open("./upload/test.zip")
	if err != nil {
		panic(err)
	}
	md5 := tool.FileMd5(open)
	fmt.Println(md5)
}

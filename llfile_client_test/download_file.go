package main

import (
	"bytes"
	"fmt"
	"io"
	"llfile_client/tool"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

//下载客户端不允许进程挂掉
func main() {
	//var reqUrl = "http://127.0.0.1:8081/file/download" //这个要鉴权
	var reqUrl = "http://127.0.0.1:8081/s/download" //这个不鉴权
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo4LCJ1c2VyX2lk55qEQUVT5re35reGIjoidFpNNnJjMFg2M0E3bHlidlF2d0dXZz09IiwiZXhwIjoxNjU5MjY4NjM0fQ.7t81Wf3Y3gB7zxC8p_vsJXhukSUawzT96WJwJd7JNi0"
	var user_id = "8"
	var eventId = "eccT9nOLJ7S9rqfC"
	var fileName = "./download/testn.zip"
	var bytess = tool.Length(fileName)

	client := &http.Client{}

	//all, err2 := io.ReadAll(reader)

	request, err2 := http.NewRequest("GET", reqUrl, nil)
	request.Header.Set("Bytes", strconv.Itoa(int(bytess)))

	request.Header.Set("Download_event_id", eventId)
	if err2 != nil {
		panic(err2)
	}

	AddUpFileHeader(request, token, user_id)

	do, err2 := client.Do(request)
	get := do.Header.Get("Error")
	if get != "" {
		fmt.Println(get)
		return
	}

	if err2 != nil {
		panic(err2)
	}

	buf := new(bytes.Buffer)
	Write(fileName, buf, do)
	io.Copy(buf, do.Body)
	fmt.Println("加载完成")
	for {

	}
}

func AddUpFileHeader(r *http.Request, token, user_id string) {
	r.Header.Set("Access-Token", token)
	r.Header.Set("User_id", user_id)
}

func Write(fileName string, buffer *bytes.Buffer, res *http.Response) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case <-c:
			//关闭连接
			fmt.Println("guan bi le lianjie")
			res.Body.Close()
			//将缓存写入文件
			file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}
			n, err := io.Copy(file, buffer)
			fmt.Println("写了多少", n)
			if err != nil {
				panic(err)
			}
			os.Exit(1)
		}
	}()
}

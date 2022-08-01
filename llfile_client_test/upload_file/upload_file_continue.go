package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var reqUrl = "http://127.0.0.1:8081/file/upload"
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo4LCJ1c2VyX2lk55qEQUVT5re35reGIjoidFpNNnJjMFg2M0E3bHlidlF2d0dXZz09IiwiZXhwIjoxNjU5MjU1MzA0fQ.xieFgyWs5iADTC9aqJ-pdC1pviSDTp9WfP9G9svdobo"
	var user_id = "8"
	var folder_id = "10"
	var eventId = "PdqvJnOhl7NzrLrA"
	var bytes = 26116096

	client := &http.Client{}
	file, err2 := os.Open("./upload/go.zip")
	file.Seek(int64(bytes), 0)
	if err2 != nil {
		panic(err2)
	}
	reader := bufio.NewReader(file)
	//all, err2 := io.ReadAll(reader)

	request, err2 := http.NewRequest("POST", reqUrl, reader)
	request.Header.Set("Context-Type", "application/octet-stream")
	request.Header.Set("Upload_event_id", eventId)
	request.Header.Set("Bytes", strconv.Itoa(bytes))
	if err2 != nil {
		panic(err2)
	}

	AddUpFileHeader2(request, token, user_id, folder_id)

	do, err2 := client.Do(request)
	if err2 != nil {
		panic(err2)
	}

	all, err2 := io.ReadAll(do.Body)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(string(all))
}

func AddUpFileHeader2(r *http.Request, token, user_id, f string) {
	r.Header.Set("Access-Token", token)
	r.Header.Set("User_id", user_id)
	r.Header.Set("Folder_id", f)
}

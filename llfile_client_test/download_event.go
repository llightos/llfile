package main

import (
	"fmt"
	"io"
	"llfile_client/tool"
	"net/http"
	"net/url"
)

func main() {
	var reqUrl = "http://127.0.0.1:8081/file/download/event"
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo4LCJ1c2VyX2lk55qEQUVT5re35reGIjoidFpNNnJjMFg2M0E3bHlidlF2d0dXZz09IiwiZXhwIjoxNjU5MjY4NjM0fQ.7t81Wf3Y3gB7zxC8p_vsJXhukSUawzT96WJwJd7JNi0"
	var user_id = "8"
	var fileName = "test.zip"

	client := &http.Client{}

	param := url.Values{}

	name, expandName := tool.DivideName(fileName)
	param.Set("head_name", name)
	param.Set("expand_name", expandName)

	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		panic(err)
	}
	request.URL.RawQuery = param.Encode()

	AddUpEventHeader(request, token, user_id)

	do, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	all, err := io.ReadAll(do.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(all))
}

func AddUpEventHeader(r *http.Request, token, user_id string) {
	r.Header.Set("Access-Token", token)
	r.Header.Set("User_id", user_id)
}

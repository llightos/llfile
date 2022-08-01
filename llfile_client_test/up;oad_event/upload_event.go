package main

import (
	"fmt"
	"io"
	"llfile_client/tool"
	"net/http"
	"net/url"
	"strconv"
)

//hash := c.Query("hash")
//		size := c.Query("size")
//		headName := c.Query("head_name")
//		expandName := c.Query("expand_name")
//		folderId := c.Query("folder_id")
//		fmt.Println("size", size)

func main() {
	var reqUrl = "http://127.0.0.1:8081/file/upload/event"
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo4LCJ1c2VyX2lk55qEQUVT5re35reGIjoidFpNNnJjMFg2M0E3bHlidlF2d0dXZz09IiwiZXhwIjoxNjU5MjU4MzA4fQ.ni9SnNPUulI_1NLcTRQeYBe69Ifipk0Zv-MyJWKs-aI"
	var user_id = "8"
	var folder_id = "10"

	client := &http.Client{}

	param := url.Values{}

	file, info := tool.OpenFile("./upload/go.zip")
	name, expandName := tool.DivideName(info.Name())
	param.Set("hash", tool.FileMd5(file))
	param.Set("head_name", name)
	param.Set("expand_name", expandName)
	param.Set("folder_id", folder_id)
	param.Set("size", strconv.Itoa(int(info.Size())))

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

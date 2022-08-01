package handler

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Context struct {
	*gin.Context
}

var UploadManager = &ManagerUpload{}

type ManagerUpload struct {
	mu sync.Mutex
	m  map[string]*UploadEvent // key是UploadEvent的id
}

type UploadEvent struct {
	Id string

	CreateTime time.Duration //事件创建时间，用来判断是否超时
	limiter    *rate.Limiter

	HeadName   string
	ExpandName string
	size       uint
	hash       string

	foldID uint

	current uint //当前完成的字节数
	io.ReadWriter

	Status uint //状态，1为正在进行，0位暂停
}

func init() {
	UploadManager = &ManagerUpload{m: make(map[string]*UploadEvent)}
}

func (u *UploadEvent) Write(p []byte) (n int, err error) { //不能让网络直接写数据，这样一旦寄，就没救了，应该设置一个缓冲
	u.limiter.Wait(context.TODO())
	n, err = u.ReadWriter.Write(p)
	if err != nil {
		log.Println(err)
	}
	u.current = u.current + uint(n)
	return
}

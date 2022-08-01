package downhandler

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"io"
	"log"
	"sync"
	"time"
)

type Context struct {
	*gin.Context
}

var DownloadManager = &ManagerDownload{}

type ManagerDownload struct {
	mu sync.Mutex //map锁

	fileMu sync.Mutex //文件读取锁

	m map[string]*DownloadEvent // key是DownloadEvent的id
}

type DownloadEvent struct {
	Id  string
	AAA bool //是否需要鉴权

	Timer   *time.Timer //事件创建时间，用来判断是否超时
	limiter *rate.Limiter

	HeadName   string
	ExpandName string
	size       uint
	hash       string

	foldID uint

	current   uint //当前完成的字节数
	io.Reader      //是bufio

	Status uint //状态，1为正在进行，0位暂停
}

func init() {
	DownloadManager = &ManagerDownload{m: make(map[string]*DownloadEvent)}
}

func (u *DownloadEvent) SetAAA(b bool) {
	u.AAA = b
}

func (u *DownloadEvent) Read(p []byte) (n int, err error) { //不能让网络直接写数据，这样一旦寄，就没救了，应该设置一个缓冲
	_ = u.limiter.Wait(context.TODO())
	n, err = u.Reader.Read(p)
	if err != nil {
		log.Println(err)
	}
	u.current = u.current + uint(n)
	return
}

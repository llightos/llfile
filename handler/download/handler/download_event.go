package downhandler

import (
	"llfile/config"

	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gogf/gf/v2/util/grand"
	"golang.org/x/time/rate"
)

// NewDownloadEvent 新建一个文件，配置这个上传事件，实现限流的配置32kb为一个单位,
// 不需要folderID
func NewDownloadEvent(headName, expandName, hash string, size uint, folderId uint) (u *DownloadEvent, err error) {
	downloadEvent := new(DownloadEvent)
	downloadEvent.HeadName = headName
	downloadEvent.ExpandName = expandName
	downloadEvent.hash = hash
	downloadEvent.size = size

	downloadEvent.Timer = time.NewTimer(20 * time.Second)
	downloadEvent.foldID = folderId

	downloadEvent.Id = grand.S(16)
	downloadEvent.limiter = rate.NewLimiter(rate.Limit(config.LimitSpeed), config.LimitSpeedInt+10) //每s生成10个->相当于320kb/s的限流
	file, err := os.Open("./file/" + hash + ".llfile")
	if err != nil {
		return nil, err
	}

	//把文件读进缓存里
	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, file)
	fmt.Println("len", buf.Len())
	_ = file.Close()
	downloadEvent.Reader = buf

	return downloadEvent, nil
}

func (u *DownloadEvent) Seek(n int) (err error) {
	defer func() {
		if pan := recover(); pan != nil {
			err = errors.New("DownloadEvent Seek Err")
		}
	}()
	// 缓存用了就没了，所以重传要重新读文件加载缓存
	open, _ := os.Open("./file/" + u.hash + ".llfile")
	seek, err := open.Seek(int64(n), 0)
	if err != nil {
		return err
	}
	newBuffer := new(bytes.Buffer)
	_, _ = io.Copy(newBuffer, open)
	_ = open.Close()
	u.Reader = newBuffer

	fmt.Println("success seek to ", seek)
	return nil
}

package handler

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/v2/util/grand"
	"golang.org/x/time/rate"
	"io"
	"llfile/config"
	"llfile/util"
	"os"
	"time"
)

// NewUploadEvent 新建一个文件，配置这个上传事件，实现限流的配置32kb为一个单位
func NewUploadEvent(headName, expandName, hash string, size uint, folderId uint) (u *UploadEvent, err error) {
	uploadEvent := new(UploadEvent)
	uploadEvent.HeadName = headName
	uploadEvent.ExpandName = expandName
	uploadEvent.hash = hash
	uploadEvent.size = size
	uploadEvent.CreateTime = time.Duration(time.Now().Unix())
	uploadEvent.foldID = folderId

	uploadEvent.Id = grand.S(16)
	uploadEvent.limiter = rate.NewLimiter(rate.Limit(config.LimitSpeed), config.LimitSpeedInt+10) //令牌桶大小为10，每s生成10个->相当于320kb/s的限流

	if size > util.GB { // 如果文件大于1Gb，就直接写硬盘
		create, err := os.Create("./file/" + u.hash + ".llfile")
		if err != nil {
			return nil, err
		}
		uploadEvent.ReadWriter = create
	} else { // 否则先写缓存再写硬盘
		b := new(bytes.Buffer)
		uploadEvent.ReadWriter = b
	}

	if err != nil {
		return nil, err
	}

	return uploadEvent, nil
}

func (u *UploadEvent) StopFile() {
	switch u.ReadWriter.(type) {
	case *bytes.Buffer:
	case *os.File:
		u.ReadWriter.(*os.File).Close()
		os.Remove("./file/" + u.hash + ".llfile")
	}
}

func (u *UploadEvent) WriteToFile() error {
	fmt.Println("是否写入？？？")
	switch u.ReadWriter.(type) {
	case *bytes.Buffer:
		create, err := os.Create("./file/" + u.hash + ".llfile")
		if err != nil {
			return err
		}
		_, err = io.Copy(create, u.ReadWriter)
		create.Close()
		if err != nil {
			return err
		}
	case *os.File:
		return nil
	}
	return nil
}

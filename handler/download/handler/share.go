package downhandler

import (
	"llfile/handler"
	"llfile/model"
	"llfile/service"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/grand"
)

type LinkInfo struct {
	Link   string //生成的链接
	Random string //链接的随机数
	Code   string
	model.File
}

// GetShareLink 获得一个文件链接
func GetShareLink() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := handler.MyCtx(context)

		userId := c.GetHeader("User_id")
		fileId := c.PostForm("file_id")
		code := c.PostForm("code")

		file := handler.DB.FindFileByFileId(userId, fileId)

		randomS := grand.S(20)
		linkInfo := new(LinkInfo)
		linkInfo.File = file
		linkInfo.Random = randomS
		linkInfo.Link = "http://127.0.0.1:8081/s?file=" + randomS
		if code != "" {
			linkInfo.Code = code
		}

		marshal, err := json.Marshal(linkInfo)
		if err != nil {
			c.ReturnErr400("json.Marshal err")
			return
		}

		err = service.RedisAdd(linkInfo.Random, string(marshal))
		if err != nil {
			c.ReturnErr400(err)
			return
		}
		c.JSON(200, gin.H{
			"status": true,
			"data":   linkInfo.Link,
			"code":   code,
		})
	}
}

func ShareDownLoadEvent() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := handler.MyCtx(context)
		random := c.Query("file")
		code := c.Query("code")

		val, ok := service.RedisVal(random)
		if !ok {
			c.ReturnErr400("no file")
			return
		}

		m := new(LinkInfo)
		_ = json.Unmarshal([]byte(val), m)

		if m.Code != code {
			c.ReturnErr("密码不正确")
		}

		event, err := NewDownloadEvent(m.Name, m.ExpandedName, m.Hash, m.Size, m.FolderID)
		event.SetAAA(false)

		if err != nil {
			c.ReturnErr400("NewDownloadEvent fail")
			return
		}

		DownloadManager.AddEvent(event)

		c.JSON(200, gin.H{
			"status":            true,
			"download_event_id": event.Id,
		})
	}
}

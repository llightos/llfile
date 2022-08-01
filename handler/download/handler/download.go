package downhandler

import (
	"llfile/handler"

	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateDownloadEvent() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := MyCtx(context)
		userId := c.GetHeader("User_id")
		fileId := c.Query("file_id")
		headName := c.Query("head_name")
		expandName := c.Query("expand_name")

		var fileHash string
		var size uint

		if fileId != "" {
			File := handler.DB.FindFileByFileId(userId, fileId)
			fileHash = File.Hash
			size = File.Size
		} else {
			File := handler.DB.FindFile(userId, headName, expandName)
			fileHash = File.Hash
			size = File.Size
		}

		if fileHash == "" {
			c.ReturnErr400("没有目标文件")
			return
		}

		event, err := NewDownloadEvent(headName, expandName, fileHash, size, 0)
		event.SetAAA(true)
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

func Download() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := handler.MyCtx(context)

		eventID := c.GetHeader("Download_event_id")

		event, ok := DownloadManager.QueryEvent(eventID)
		//fmt.Println("ok?", ok)
		if !ok {
			fmt.Println("!ok!!")
			c.ReturnErrHeader("DownloadManager.QueryEvent寄")
			c.ReturnErr400("DownloadManager.QueryEvent寄")
			c.Abort()
			return
		}

		bytes := c.GetHeader("Bytes")

		//说明需要鉴权
		if event.AAA == true {
			_, exists := c.Get("User_id")
			if !exists {
				c.ReturnErr400("非法接口")
				return
			}
		}

		//说明是断的传
		if bytes != "" {
			atoi, err := strconv.Atoi(bytes)
			if err != nil {
				c.ReturnErr(err)
				return
			}
			err = event.Seek(atoi)
			if err != nil {
				c.ReturnErr(err)
				return
			}
		}

		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-MD5", event.hash)
		c.Writer.Header().Set("Content-Length", strconv.Itoa(int(event.size)))

		written, err := io.Copy(c.Writer, event)
		fmt.Println("written", written)
		if err != nil {
			log.Println("断开了连接")
			return
		}
		DownloadManager.DeleteEvent(eventID)
	}
}

//func GetUploadBytes() {
//
//}

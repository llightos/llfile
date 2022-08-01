package handler

import (
	"llfile/handler"
	"llfile/model"
	"llfile/msg"
	"llfile/util"

	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUploadEvent() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := handler.MyCtx(context)
		hash := c.Query("hash")
		size := c.Query("size")
		headName := c.Query("head_name")
		expandName := c.Query("expand_name")
		folderId := c.Query("folder_id")
		fmt.Println("size", size)
		newFileSize, err := strconv.Atoi(size)
		if err != nil {
			c.ReturnErr("param not right")
			return
		}
		bool1, bool2 := model.NewModelDB().IfHashExist(c.GetHeader("User_id"), hash)

		// 用户已经有这个文件了，
		if bool2 {
			c.ReturnErr("你已经有这个文件了")
			return
		}

		folderID, err := strconv.Atoi(folderId)
		if err != nil {
			c.ReturnErr("param not right")
			return
		}

		//已经存在hash值的情况
		if bool1 {
			fileId, err := model.NewModelDB().AddFile(model.File{
				FolderID:     uint(folderID),
				Hash:         hash,
				Name:         headName,
				ExpandedName: expandName,
				Size:         uint(newFileSize),
				Path:         model.FolderPath(uint(folderID)),
			})
			if err != nil {
				c.ReturnErr(err)
				return
			}
			c.JSON(200, gin.H{
				"status": true,
				"fileId": fileId,
			})
			return
		}

		uploadEvent, err := NewUploadEvent(headName, expandName, hash, uint(newFileSize), uint(folderID))
		if err != nil {
			c.ReturnErr("create file false")
			return
		}
		UploadManager.AddEvent(uploadEvent)
		c.JSON(200, gin.H{
			"status":          true,
			"upload_event_id": uploadEvent.Id,
		})
	}
}

func UpdateContinue() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := handler.MyCtx(context)
		_, exists := c.Get("do")
		if !exists {
			return
		}
		contextType := c.GetHeader("Context-Type")
		eventID := c.GetHeader("Upload_event_id")
		if contextType != "application/octet-stream" && contextType != "binary" {
			c.ReturnErr("Context-Type err")
			return
		}
		event, ok := UploadManager.QueryEvent(eventID)
		if !ok {
			c.ReturnErr400("没有这个上传事件")
			return
		}

		bytes := c.GetHeader("Bytes")
		fmt.Println("event.bytes, bytes", event.current, bytes)
		if bytes != strconv.Itoa(int(event.current)) {
			c.ReturnErr(msg.ErrUploadEventId)
			return
		}

		//开始上传
		_, err := io.Copy(event, c.Request.Body)
		if err != nil {
			c.ReturnErr(msg.ErrUploadFile)
			return
		}

		err = event.WriteToFile() //正式写入服务器硬盘

		if !util.FileHash(event.hash) { //计算存储后的文件与传过来的hash不相等
			//如果不晓得就删除文件
			fmt.Println("hash:", event.hash)
			_ = os.Remove("./file/" + event.hash + "./llfile")
			event.StopFile()
			c.ReturnErr400("文件hash值与入参不匹配")
			return
		}

		file, err := model.NewModelDB().AddFile(model.File{
			FolderID:     event.foldID,
			Hash:         event.hash,
			Name:         event.HeadName,
			ExpandedName: event.ExpandName,
			Size:         event.size,
			Path:         model.FolderPath(event.foldID),
		})

		if err != nil {
			c.ReturnErr("寄")
			return
		}

		c.JSON(200, gin.H{
			"status": true,
			"data":   file,
		})
		//event.StopFile()
		UploadManager.DeleteEvent(eventID)
	}
}

func Upload() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := handler.MyCtx(context)
		contextType := c.GetHeader("Context-Type")
		//fmt.Println("header", contextType)
		eventID := c.GetHeader("Upload_event_id")
		folder := c.GetHeader("Folder_id")
		bytes := c.GetHeader("Bytes")
		if contextType != "application/octet-stream" && contextType != "binary" {
			c.ReturnErr("Context-Type err")
			return
		}

		uploadEvent, ok := UploadManager.QueryEvent(eventID)
		if !ok {
			c.ReturnErr("no exist event")
			return
		}

		if uploadEvent.current != 0 { //此时表示文件已经存储了字节了，要求使用其他接口进行上传续传
			if bytes == "" {
				c.JSON(200, gin.H{
					"WARNING": "请使用其他接口",
					"bytes":   uploadEvent.current,
				})
			}
			c.Set("do", "ok")
			c.Next()
			return
		}

		_, err := io.Copy(uploadEvent, c.Request.Body)
		if err != nil {
			c.ReturnErr500("io.Copy err")
			return
		}

		folderId, err := strconv.Atoi(folder)
		if err != nil {
			c.ReturnErr400(err)
			return
		}

		err = uploadEvent.WriteToFile() //正式写入服务器硬盘
		if err != nil {
			c.ReturnErr500("服务器寄了")
			return
		}

		if !util.FileHash(uploadEvent.hash) { //计算存储后的文件与传过来的hash不相等
			//如果不晓得就删除文件
			_ = os.Remove("./file/" + uploadEvent.hash + "./llfile")
			c.ReturnErr400("文件hash值与入参不匹配")
			return
		}

		fileid, err := model.NewModelDB().AddFile(model.File{
			FolderID:     uint(folderId),
			Hash:         uploadEvent.hash,
			Name:         uploadEvent.HeadName,
			ExpandedName: uploadEvent.ExpandName,
			Size:         uploadEvent.size,
			Path:         model.FolderPath(uint(folderId)),
		})

		if err != nil {
			c.ReturnErr400(err)
			return
		}

		c.JSON(200, gin.H{
			"status":  true,
			"file_id": fileid,
		})
		UploadManager.DeleteEvent(eventID)
		c.Abort()
		return
	}
}

func GetUploadBytes() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := MyCtx(context)
		uploadEventId := c.Query("Upload_event_id")
		event, ok := UploadManager.QueryEvent(uploadEventId)
		if !ok {
			c.ReturnErr("no Upload_event_id")
			return
		}
		c.JSON(200, gin.H{
			"status": true,
			"bytes":  event.current,
		})
	}
}

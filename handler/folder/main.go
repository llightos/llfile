package folder

import (
	"llfile/handler"
	"llfile/model"

	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FoldersAndFiles struct {
	Files   []model.QueryUserFiles
	Folders []model.QueryUserFolders
}

func ReturnFoldersAndFiles() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := handler.MyCtx(context)
		//userId, _ := c.Get("User_id")

		files := handler.DB.QueryUserFiles(c.GetHeader("User_id"))
		folders := handler.DB.QueryUserFolders(c.GetHeader("User_id"))
		f := new(FoldersAndFiles)
		f.Files = files
		f.Folders = folders
		c.JSON(200, gin.H{
			"data": f,
		})
	}
}

func AddFolder() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := handler.MyCtx(context)
		UserId := c.GetHeader("User_id")

		folderId := c.Query("folder_id")
		newName := c.Query("new_folder_name")

		if newName != "" {
			c.ReturnErr("文件夹名字不能为空")
			return
		}

		FID, err := strconv.Atoi(folderId)
		if err != nil {
			c.ReturnErr400(err)
		}
		userIDInt, err := strconv.Atoi(UserId)
		if err != nil {
			c.ReturnErr400(err)
			return
		}

		_, err = handler.DB.AddChild(uint(FID), uint(userIDInt), newName)
		if err != nil {
			c.ReturnErr400(err)
			return
		}

		c.JSON(200, gin.H{
			"status": true,
			"data":   "success",
		})

	}
}

//func ChanFileName() gin.HandlerFunc {
//	return func(context *gin.Context) {
//		c := handler.MyCtx(context)
//		c.
//	}
//}

func ChangeFileRoute() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := handler.MyCtx(context)

		get, _ := c.Get("User_id")
		fmt.Println("get？", get)
		fileId := c.PostForm("file_id")
		targetFolderId := c.PostForm("target")
		newHead := c.PostForm("new_head")
		newExpand := c.PostForm("new_expand")

		fileIdInt, err := strconv.Atoi(fileId)
		if err != nil {
			c.ReturnErr400(err)
			return
		}

		var targetFolderIdInt int

		if targetFolderId != "" {
			targetFolderIdInt, err = strconv.Atoi(targetFolderId)
			if err != nil {
				c.ReturnErr400(err)
				return
			}
		}

		err = handler.DB.ChangeFileRoute(fileIdInt, targetFolderIdInt, get.(int), newHead, newExpand)
		if err != nil {
			fmt.Println(err)
			c.ReturnErr400(err)
			return
		}
		c.JSON(200, gin.H{
			"status": true,
			"data":   "success",
		})
	}
}

//func AddFolder() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		model.NewModelDB().AddChild()
//	}
//}

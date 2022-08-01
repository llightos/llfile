package cmd

import (
	"llfile/handler"
	downhandler "llfile/handler/download/handler"
	"llfile/handler/folder"
	uphandler "llfile/handler/upload/handler"

	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()

	//登录注册
	userGroup := r.Group("user")
	{
		userGroup.POST("register", handler.Register())
		userGroup.POST("login", handler.Login())
	}

	//文件上传
	uploadGroup := r.Group("file/upload")
	uploadGroup.Use(handler.TokenVerify())
	{
		uploadGroup.GET("event", uphandler.CreateUploadEvent())
		uploadGroup.POST("", uphandler.Upload(), uphandler.UpdateContinue())
		uploadGroup.GET("event/bytes", uphandler.GetUploadBytes())
	}

	//文件下载
	downloadGroup := r.Group("file/download")
	downloadGroup.Use(handler.TokenVerify())
	{
		downloadGroup.GET("event", downhandler.CreateDownloadEvent())
		downloadGroup.GET("", downhandler.Download())
	}

	//文件管理
	queryGroup := r.Group("folder")
	queryGroup.Use(handler.TokenVerify())
	{
		queryGroup.GET("", folder.ReturnFoldersAndFiles())
		queryGroup.GET("add", folder.AddFolder())
		queryGroup.POST("file/change", folder.ChangeFileRoute())
	}

	r.POST("/share/link", handler.TokenVerify(), downhandler.GetShareLink())
	//不鉴权的获取
	r.GET("/s", downhandler.ShareDownLoadEvent())
	//不鉴权的下载
	r.GET("/s/download", downhandler.Download())

	r.Run(":8081")
}

package handler

import (
	"github.com/gin-gonic/gin"
	"llfile/rpc"
)

func Register() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := MyCtx(context)
		userName := c.PostForm("username")
		passWord := c.PostForm("password")

		ok, id := rpc.NewUser().CallRegister(userName, passWord)
		if !ok {
			c.ReturnErr("注册失败")
			return
		}
		c.JSON(200, gin.H{
			"status":  true,
			"user_id": id,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := MyCtx(context)
		username := c.PostForm("username")
		password := c.PostForm("password")
		ok, token, u := rpc.NewUser().CallLogin(username, password)
		if !ok {
			c.ReturnErr400("登录")
			return
		}
		c.JSON(200, gin.H{
			"status": true,
			"token":  token,
			"userid": u,
		})
	}
}

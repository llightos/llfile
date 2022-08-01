package handler

import (
	"github.com/gin-gonic/gin"
)

func MyCtx(c *gin.Context) *Context {
	ctx := new(Context)
	ctx.Context = c
	return ctx
}

// ReturnErr 返回标准的错误信息
func (c *Context) ReturnErr(data string) {
	c.JSON(200, gin.H{
		"status": false,
		"data":   data,
	})
	c.Abort()
}

func (c *Context) ReturnErr500(data string) {
	c.JSON(500, gin.H{
		"status": false,
		"data":   data,
	})
	c.Abort()
}

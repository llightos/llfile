package downhandler

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func MyCtx(c *gin.Context) *Context {
	ctx := new(Context)
	ctx.Context = c
	return ctx
}

// ReturnErr 返回标准的错误信息
func (c *Context) ReturnErr(data interface{}) {
	c.JSON(200, gin.H{
		"status": false,
		"data":   data,
	})
}

func (c *Context) ReturnErr500(data string) {
	c.JSON(500, gin.H{
		"status": false,
		"data":   data,
	})
}

func (c *Context) ReturnErr400(data interface{}) {
	c.JSON(400, gin.H{
		"status": false,
		"data":   data,
	})
}

// String2Uint 将string转换为uint，转换失败就返回
func (c *Context) String2Uint(s string) uint {
	atoi, err := strconv.Atoi(s)
	if err != nil {
		c.JSON(401, gin.H{
			"status": false,
			"data":   "参数格式错误",
		})
		c.Abort()
	}
	return uint(atoi)
}

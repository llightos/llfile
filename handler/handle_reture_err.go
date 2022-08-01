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
func (c *Context) ReturnErr(data interface{}) {
	c.JSON(200, gin.H{
		"status": false,
		"data":   data,
	})
	c.Abort()
}

func (c *Context) ReturnErr400(data interface{}) {
	c.JSON(400, gin.H{
		"status": false,
		"data":   data,
	})
	c.Abort()
}

func (c *Context) ReturnErrHeader(data string) {
	c.Writer.Header().Add("Error", data)
}

func (c *Context) ReturnErr500(data string) {
	c.JSON(500, gin.H{
		"status": false,
		"data":   data,
	})
	c.Abort()
}

//func (c *Context)String2Uint(data string) uint {
//	atoi, err := strconv.Atoi(data)
//	if err != nil{
//		c.JSON(400, gin.H{
//			"status" : false,
//			"data" : "param not right",
//		})
//		c.Abort()
//		return 0
//	}
//	return uint(atoi)
//}

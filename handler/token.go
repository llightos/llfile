package handler

import (
	"github.com/gin-gonic/gin"
	"llfile/rpc"
	"strconv"
)

type Context struct {
	*gin.Context
}

func TokenVerify() gin.HandlerFunc {
	return func(context *gin.Context) {
		c := MyCtx(context)
		token := c.GetHeader("Access-Token")
		userID := c.GetHeader("User_id")
		//fmt.Println("token, user_id", token, userID)

		atoi, err := strconv.Atoi(userID)
		if err != nil {
			c.ReturnErr("err user_id")
			c.Abort()
			return
		}
		authToken, err := rpc.NewAuth().CallAuthToken(token, atoi)
		if err != nil {
			c.ReturnErr400(authToken.Data)
			c.Abort()
			return
		}
		//fmt.Println("ok??????????", authToken.Ok)
		if authToken.Ok == false {
			c.ReturnErr400(authToken.Data)
			c.Abort()
			return
		}
		c.Set("User_id", atoi)
	}
}

package middleware

import (
	"GoProject/global"
	"GoProject/util"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(global.Logger.Out, defaultHandleRecovery)
}

func defaultHandleRecovery(c *gin.Context, err any) {
	//c.JSON(util.UnKnowResponse(err.(string)))
	c.JSON(util.UnKnowResponse(""))
}

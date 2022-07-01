package article

import (
	"github.com/gin-gonic/gin"
)

func InitRoute(rg *gin.RouterGroup) {
	userRouter := rg.Group("/article")

	userRouter.GET("", defaultHandler)
	userRouter.GET("/overflow", defaultHandler)
	userRouter.GET("/list", defaultHandler)
}

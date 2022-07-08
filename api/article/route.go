package article

import (
	"github.com/gin-gonic/gin"
)

func InitRoute(rg *gin.RouterGroup) {
	articleRouter := rg.Group("/article")

	articleRouter.GET("", defaultHandler)
	articleRouter.GET("/overflow", defaultHandler)
	articleRouter.GET("/list", defaultHandler)
}

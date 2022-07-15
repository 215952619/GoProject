package article

import (
	"GoProject/middleware"
	"GoProject/util"
	"github.com/gin-gonic/gin"
	"sync"
)

var localLock sync.Locker

func InitRoute(rg *gin.RouterGroup) {
	articleRouter := rg.Group("/article")

	articleRouter.GET("", defaultHandler)
	articleRouter.GET("/overflow", util.ResponseWarp(overflow))
	articleRouter.GET("/list", util.ResponseWarp(list))
	articleRouter.GET("/article/:id", util.ResponseWarp(detail))

	articleRouter.POST("/create", util.ResponseWarp(create))
	articleRouter.Use(middleware.LogonOnly())
	articleRouter.PUT("/article/:id/recommend", defaultHandler)
	articleRouter.PUT("/article/:id/unrecommended", defaultHandler)

	articleRouter.Use(middleware.AdminOnly())
	articleRouter.PUT("/article/:id/top", defaultHandler)
	articleRouter.PUT("/article/:id/unto", defaultHandler)
}

package article

import (
	"GoProject/middleware"
	"GoProject/util"
	"github.com/gin-gonic/gin"
)

func InitRoute(rg *gin.RouterGroup) {
	articleRouter := rg.Group("/article")

	articleRouter.GET("/overflow", util.ResponseWarp(overflow))
	articleRouter.GET("/list", util.ResponseWarp(list))
	articleRouter.GET("/article/:id", util.ResponseWarp(detail))

	articleRouter.Use(middleware.LogonOnly())
	articleRouter.POST("/create", util.ResponseWarp(create))
	articleRouter.PUT("/article/:id/recommend", util.ResponseWarp(recommend))
	articleRouter.PUT("/article/:id/unrecommended", util.ResponseWarp(unrecommended))

	articleRouter.Use(middleware.AdminOnly())
	articleRouter.PUT("/article/:id/top", util.ResponseWarp(top))
	articleRouter.PUT("/article/:id/untop", util.ResponseWarp(untop))
}

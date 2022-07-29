package article

import (
	"GoProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute(rg *gin.RouterGroup) {
	articleRouter := rg.Group("/article")

	articleRouter.GET("/overflow", middleware.ResponseWarp(overflow))
	articleRouter.GET("/list", middleware.ResponseWarp(list))
	articleRouter.GET("/article/:id", middleware.ResponseWarp(detail))

	articleRouter.Use(middleware.LogonOnly())
	articleRouter.POST("/create", middleware.ResponseWarp(create))
	articleRouter.PUT("/article/:id/recommend", middleware.ResponseWarp(recommend))
	articleRouter.PUT("/article/:id/unrecommended", middleware.ResponseWarp(unrecommended))

	articleRouter.Use(middleware.AdminOnly())
	articleRouter.PUT("/article/:id/top", middleware.ResponseWarp(top))
	articleRouter.PUT("/article/:id/untop", middleware.ResponseWarp(untop))
}

package user

import (
	"GoProject/middleware"
	"GoProject/util"
	"github.com/gin-gonic/gin"
)

func InitRoute(rg *gin.RouterGroup) {
	userRouter := rg.Group("/user")

	userRouter.POST("/logon", util.ResponseWarp(logon))
	userRouter.GET("/sso/:platform", util.ResponseWarp(getCode))
	userRouter.GET("/sso/:platform/redirect", util.ResponseWarp(ssoRedirect))

	userRouter.Use(middleware.LogonOnly())

	userRouter.POST("/create", util.ResponseWarp(createUser))
	userRouter.GET("/list", middleware.AdminOnly(), util.ResponseWarp(userList))
	userRouter.GET("/:id", middleware.AdminOnly(), util.ResponseWarp(userDetail))
}

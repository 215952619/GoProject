package user

import (
	"GoProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute(rg *gin.RouterGroup) {
	userRouter := rg.Group("/user")
	userRouter.POST("/login", middleware.ResponseWarp(login))
	userRouter.GET("/sso/:platform", middleware.ResponseWarp(getCode))
	userRouter.GET("/sso/:platform/redirect", middleware.ResponseWarp(ssoRedirect))

	userRouter.Use(middleware.LogonOnly())
	userRouter.POST("/create", middleware.ResponseWarp(createUser))
	userRouter.GET("/list", middleware.AdminOnly(), middleware.ResponseWarp(userList))

	userRouter.Use(middleware.AdminOnly())
	userRouter.GET("/:id", middleware.ResponseWarp(userDetail))
}

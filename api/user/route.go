package user

import (
	"GoProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute(rg *gin.RouterGroup) {
	userRouter := rg.Group("/user")
	userRouter.Use(middleware.LogonOnly())

	userRouter.GET("", defaultHandler)
	userRouter.GET("/list", middleware.AdminOnly(), userList)
	userRouter.GET("/:id", middleware.AdminOnly(), userDetail)
}

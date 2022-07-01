package api

import (
	"GoProject/api/user"
	"GoProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) {
	apiRouter := r.Group("/api")
	apiRouter.Use(middleware.Resolve())
	apiRouter.GET("/user/logon", logon)
	user.InitRoute(apiRouter)
}

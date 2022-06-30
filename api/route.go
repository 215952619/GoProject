package api

import (
	"GoProject/api/user"
	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) {
	apiRouter := r.Group("/api")
	apiRouter.GET("/user/logon", logon)
	user.InitRoute(apiRouter)
}

package api

import (
	"GoProject/api/article"
	"GoProject/api/user"
	"GoProject/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.Engine) {
	apiRouter := r.Group("/api")
	apiRouter.Use(middleware.Resolve())
	user.InitRoute(apiRouter)
	article.InitRoute(apiRouter)
}

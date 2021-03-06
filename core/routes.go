package core

import (
	"GoProject/api"
	"GoProject/global"
	"GoProject/middleware"
	"embed"
	"github.com/gin-gonic/gin"
)

func InitRoutes(sources *embed.FS) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.Cors())
	r.Static("/static", "./static")
	api.InitRoute(r)

	if global.Mode == "product" {
		r.Use(middleware.HtmlRender("/", *sources))
	}

	return r
}

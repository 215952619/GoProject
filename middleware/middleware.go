package middleware

import (
	"GoProject/global"
	"GoProject/util"
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！） TODO 需要区分环境
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,Content-Type")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}
		c.Next()
	}
}

func HtmlRender(urlPrefix string, fs embed.FS) gin.HandlerFunc {
	const indexHtml = "index.html"
	const basePath = "frontend/dist"

	return func(c *gin.Context) {
		urlPath := strings.TrimSpace(c.Request.URL.Path)
		if urlPath == urlPrefix {
			urlPath = path.Join(urlPrefix, indexHtml)
		}
		urlPath = filepath.Join(basePath, urlPath)

		if os.Getenv("GOOS") == "windows" {
			urlPath = filepath.ToSlash(urlPath) // TODO: 针对Windows系统
		}

		f, err := fs.Open(urlPath)
		if err != nil {
			return
		}
		fi, err := f.Stat()
		if strings.HasSuffix(urlPath, ".html") {
			c.Header("Cache-Control", "no-cache")
			c.Header("Content-Type", "text/html; charset=utf-8")
		}

		if strings.HasSuffix(urlPath, ".js") {
			c.Header("Content-Type", "text/javascript; charset=utf-8")
		}
		if strings.HasSuffix(urlPath, ".css") {
			c.Header("Content-Type", "text/css; charset=utf-8")
		}

		if err != nil || !fi.IsDir() {
			bs, err := fs.ReadFile(urlPath)
			if err != nil {
				global.Logger.WithFields(logrus.Fields{
					"url": urlPath,
				}).Panic("")
			}
			c.Status(200)
			c.Writer.Write(bs)
			c.Abort()
		}
	}
}

func LogonOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *global.User
		user, err := util.ParseJwt(c)
		if err != nil {
			//todo sign
			c.JSON(http.StatusOK, gin.H{"msg": "not match"})
			c.Abort()
		} else {
			if err := global.DB.Where("id=?", user.ID).First(&user).Error; err != nil {
				//not found
				c.JSON(http.StatusOK, nil)
				c.Abort()
			}
			if user.Status == global.Closed {
				//in black list
				c.JSON(http.StatusOK, nil)
				c.Abort()
			}

			c.Set(global.AuthedKey, user)
			c.Next()
		}
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, exists := c.Get(global.AuthedKey)
		if exists {
			user := data.(global.User)
			if user.Role != global.Normal {
				//permission denied
				c.JSON(http.StatusOK, nil)
				c.Abort()
			}
			c.Next()
		}
		// not logon
		c.JSON(http.StatusOK, nil)
		c.Abort()
	}
}

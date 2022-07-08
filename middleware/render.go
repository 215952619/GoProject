package middleware

import (
	"GoProject/global"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func HtmlRender(urlPrefix string, fs embed.FS) gin.HandlerFunc {
	const indexHtml = "index.html"
	const basePath = "frontend/dist"

	return func(c *gin.Context) {

		fmt.Println("render")
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
			global.Logger.WithFields(logrus.Fields{
				"err":       err.Error(),
				"urlPrefix": urlPrefix,
				"urlPath":   urlPath,
			}).Error("open static file error")
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
					"err":       err.Error(),
					"urlPrefix": urlPrefix,
					"urlPath":   urlPath,
				}).Error("read static file error")
			}
			c.Status(200)
			c.Writer.Write(bs)
			c.Abort()
		}
	}
}

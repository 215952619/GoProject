package middleware

import (
	"GoProject/database"
	"GoProject/global"
	"GoProject/util"
	"github.com/gin-gonic/gin"
)

func LogonOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get(global.AuthedKey)
		if exists {
			c.Next()
		} else {
			// not logon
			c.JSON(util.DefaultResponse(util.NotLogon))
			c.Abort()
		}
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, exists := c.Get(global.AuthedKey)
		if exists {
			user := data.(database.User)
			if user.Role == database.Normal {
				//permission denied
				c.JSON(util.DefaultResponse(util.PermissionDenied))
				c.Abort()
			} else {
				c.Next()
			}
		} else {
			// not logon
			c.JSON(util.DefaultResponse(util.NotLogon))
			c.Abort()
		}
	}
}

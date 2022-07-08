package middleware

import (
	"GoProject/database"
	"GoProject/global"
	"GoProject/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

func LogonOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("logon only")
		_, exists := c.Get(global.AuthedKey)
		if exists {
			c.Next()
		} else {
			// not logon
			c.JSON(util.CustomResponse(util.NotLogonErrorTemplate, "", nil))
			c.Abort()
		}
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("admin only")
		data, exists := c.Get(global.AuthedKey)
		if exists {
			user := data.(database.User)
			if user.Role == database.Normal {
				//permission denied
				c.JSON(util.CustomResponse(util.PermissionDeniedErrorTemplate, "", nil))
				c.Abort()
			} else {
				c.Next()
			}
		} else {
			// not logon
			c.JSON(util.CustomResponse(util.NotLogonErrorTemplate, "", nil))
			c.Abort()
		}
	}
}

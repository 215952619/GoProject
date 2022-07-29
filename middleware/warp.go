package middleware

import (
	"GoProject/database"
	"GoProject/global"
	"GoProject/util"
	"github.com/gin-gonic/gin"
)

type routeHandler func(c *gin.Context, user *database.User) (data interface{}, err error)

func ResponseWarp(cb routeHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user database.User
		if data, exists := c.Get(global.AuthedKey); exists {
			user = data.(database.User)
		}
		data, err := cb(c, &user)

		if err != nil {
			errorCode := util.ErrorToErrorString(err)
			c.JSON(util.CustomResponse(errorCode, "", nil))
			return
		}

		c.JSON(util.SuccessResponse(data))
	}
}

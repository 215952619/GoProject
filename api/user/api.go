package user

import (
	"GoProject/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "not Implementation"})
}

func userList(c *gin.Context) {
	var users []global.User
	if global.DBM.List(&users) != nil {
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusOK, users)
	}
}

func userDetail(c *gin.Context) {
	id := c.Param("id")
	if len(id) > 0 {
		var user global.User
		if err := global.DBM.First(&user, "id=?", id); err != nil {
			//	not match
			c.JSON(http.StatusOK, nil)
		} else {
			c.JSON(http.StatusOK, user)
		}
	} else {
		//params error
		c.JSON(http.StatusOK, nil)
	}
}

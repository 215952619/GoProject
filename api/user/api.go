package user

import (
	"GoProject/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}

func userList(c *gin.Context) {
	var users []global.User
	if err := global.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{})
	}
	c.JSON(http.StatusOK, users)
}

func userDetail(c *gin.Context) {

}

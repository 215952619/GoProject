package user

import (
	"GoProject/database"
	"GoProject/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "not Implementation"})
}

func createUser(c *gin.Context) {
	var newUser NewUserRequest
	if err := c.ShouldBind(&newUser); err != nil {
		if err := database.DBM.Create(&newUser); err != nil {
			c.JSON(util.UnKnowResponse(err.Error()))
		} else {
			c.JSON(util.SuccessResponse(nil))
		}
	}
}

func userList(c *gin.Context) {
	var users []database.User
	if database.DBM.List(&users) != nil {
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusOK, users)
	}
}

func userDetail(c *gin.Context) {
	id := c.Param("id")
	if len(id) > 0 {
		var user database.User
		if err := database.DBM.First(&user, "id=?", id); err != nil {
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

package api

import "github.com/gin-gonic/gin"

func logon(c *gin.Context) {
	var form logonRequest
	c.ShouldBind(&form)
	form.CheckCode()
}

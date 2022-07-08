package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func defaultHandler(c *gin.Context) {
	panic("panic msg")
	c.JSON(http.StatusOK, gin.H{"msg": "not Implementation"})
}

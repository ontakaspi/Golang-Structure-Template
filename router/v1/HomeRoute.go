package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetHomeRoutes(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "This is curent version 1 active path url "})
	})
}

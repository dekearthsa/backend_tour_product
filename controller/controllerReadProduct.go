package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ControllerReadProduct(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ControllerDebug(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "service tour product is ok."})
}

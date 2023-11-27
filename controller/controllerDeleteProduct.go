package controller

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

type modelKey struct {
	ProductName string
}

func ControllerDeleteProduct(c *gin.Context) {
	const PROJECTID = "confident-topic-404213"
	const KIND = "product"

	ctx := context.Background()
	var req modelKey
	if err := c.BindJSON(&req); err != nil {
		log.Println("err BindJSON =>", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad request."})
	}

	client, err := datastore.NewClient(ctx, PROJECTID)
	if err != nil {
		log.Println("err create client datastore", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "err create client datastore."})
	}
	defer client.Close()

	keyEntity := datastore.NameKey(KIND, req.ProductName, nil)
	if err := client.Delete(ctx, keyEntity); err != nil {
		log.Println("cannot delete => ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"stataus": "cannot delete from server. "})
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

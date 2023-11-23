package controller

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"project_shopping_tour/service_product/model"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

func ControllerReadProduct(c *gin.Context) {
	const PROJECTID = "confident-topic-404213"
	const KIND = "product"
	const BUCKET = "padtravel"
	ctx := context.Background()
	clientDatastore, err := datastore.NewClient(ctx, PROJECTID)
	if err != nil {
		log.Println(err)
	}
	var data []model.Products
	var resultOut []model.ProductOut

	keys, err := clientDatastore.GetAll(ctx, datastore.NewQuery(KIND), &data)
	if err != nil {
		log.Println("err in ControllerGetAllPath => ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": err})
	}

	for i, key := range keys {
		log.Println(key)
		log.Println(i)
	}
	for i, element := range data {

		var arrayBase64 []string

		for _, item := range element.ImagePath {
			client, err := storage.NewClient(ctx)
			if err != nil {
				log.Println("can't create new client ", err)
			}
			buckets := client.Bucket(BUCKET)
			rc, err := buckets.Object(item).NewReader(ctx)
			if err != nil {
				log.Println("err when fetch image from bucket", err)
				c.JSON(http.StatusServiceUnavailable, gin.H{"Status": "err when fetch image from bucket."})
			}
			byteFile, err := io.ReadAll(rc)
			defer rc.Close()
			if err != nil {
				log.Println("err read file from bucket")
				c.JSON(http.StatusInternalServerError, gin.H{"Status": "err read file from bucket."})
			}
			str := base64.StdEncoding.EncodeToString(byteFile)
			arrayBase64 = append(arrayBase64, str)
		}

		dataOut := model.ProductOut{
			Region:       data[i].Region,
			ProductName:  data[i].ProductName,
			Objective:    data[i].Objective,
			Introduction: data[i].Introduction,
			Price:        data[i].Price,
			Include:      data[i].Include,
			Exclusive:    data[i].Exclusive,
			Content:      data[i].Content,
			ImagePath:    data[i].ImagePath,
			ArrayBase64:  arrayBase64,
		}

		resultOut = append(resultOut, dataOut)

	}

	c.JSON(http.StatusOK, gin.H{"reply": resultOut})
}

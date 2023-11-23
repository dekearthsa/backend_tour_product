package controller

import (
	"context"
	"io"
	"log"
	"net/http"
	"project_shopping_tour/service_product/model"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

func ControllerUpdateProduct(c *gin.Context) {
	const PROJECTID = "confident-topic-404213"
	const KIND = "content"
	const BUCKET = "padtravel"

	var timing = time.Now().UnixNano()
	var arrayImagePath []string
	var setPrice []model.PricePerPerson
	var setContent []model.ContentDate
	ctx := context.Background()
	// var setTitle string
	// var setContent string
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("err MultipartForm => ", err)
	}

	files := form.File["images"]
	region := form.Value["region"]
	productName := form.Value["productName"]
	objective := form.Value["objective"]
	introduction := form.Value["introduction"]
	include := form.Value["include"]
	exclusive := form.Value["exclusive"]
	persons := form.Value["person"]
	prices := form.Value["price"]
	titles := form.Value["title"]
	contents := form.Value["content"]

	for _, file := range files {
		size := file.Size
		if size >= 5000000 {
			log.Println("error file to big.")
			c.JSON(http.StatusRequestHeaderFieldsTooLarge, gin.H{"Status": "file must less than 5MB."})
		}
		src, err := file.Open()
		if err != nil {
			log.Println("err open file => ", err)
		}
		defer src.Close()
		imagePath := productName[0] + "_" + file.Filename + strconv.Itoa(int(timing))
		arrayImagePath = append(arrayImagePath, imagePath)
		clientStorage, err := storage.NewClient(ctx)
		if err != nil {
			log.Print("err in create client cloud storage => ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "err in create client cloud storage"})
		}
		bucket := clientStorage.Bucket(BUCKET)
		wc := bucket.Object(imagePath).NewWriter(ctx)
		_, err = io.Copy(wc, src)
		if err != nil {
			log.Println("err when write in cloud storage => ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "err when write in cloud storage"})
		}
		err = wc.Close()
		if err != nil {
			log.Println("err when close cloud storage => ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "err when close cloud storage"})
		}
	}

	for idx, person := range persons {
		var isPrice = prices[idx]
		createStruct := model.PricePerPerson{
			Person: person,
			Price:  isPrice,
		}
		// log.Println(createStruct)
		setPrice = append(setPrice, createStruct)
	}

	for idx, title := range titles {
		var isContent = contents[idx]
		createStruct := model.ContentDate{
			Title:   title,
			Content: isContent,
		}
		// log.Println(createStruct)
		setContent = append(setContent, createStruct)
	}

	payload := model.Products{
		Region:       region[0],
		ProductName:  productName[0],
		Objective:    objective,
		Introduction: introduction[0],
		Price:        setPrice,
		Include:      include,
		Exclusive:    exclusive,
		Content:      setContent,
		ImagePath:    arrayImagePath,
	}

	// log.Print("payload => ", payload)

	client, err := datastore.NewClient(ctx, PROJECTID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "cannot connet datastore clinet"})
	}

	keyEntity := datastore.NameKey(KIND, productName[0], nil)
	tx, err := client.NewTransaction(ctx)
	if err != nil {
		log.Println("client.NewTransaction => ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "NewTransaction error."})
	}

	if _, err := tx.Put(keyEntity, &payload); err != nil {
		log.Println("tx.Put => ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "put serror."})
	}

	if _, err := tx.Commit(); err != nil {
		log.Println("tx.Commit => ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Commit error."})
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

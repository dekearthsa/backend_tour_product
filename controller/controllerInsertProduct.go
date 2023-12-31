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

func ControllerInsertProduct(c *gin.Context) {
	const PROJECTID = "confident-topic-404213"
	const KIND = "product"
	const BUCKET = "padtravel"

	var timing = time.Now().UnixNano()
	var arrayImgPath []string
	var setPrice []model.PricePerPerson
	var setContent []model.ContentDate

	ctx := context.Background()
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
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
	// log.Println(contents)

	clientDatastore, err := datastore.NewClient(ctx, PROJECTID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "can't find projectID."})
	}

	// var setKey = strings.Join(titles[0] , "")
	key := datastore.NameKey(KIND, productName[0], nil)

	for _, file := range files {
		size := file.Size
		if size >= 5000000 {
			log.Println("error file to big.")
			c.JSON(http.StatusRequestHeaderFieldsTooLarge, gin.H{"Status": "file must less than 5MB."})
		}

		src, err := file.Open()
		if err != nil {
			log.Println("err can't open file", err)
			c.JSON(http.StatusInternalServerError, gin.H{"Status": "file can't open"})
		}
		defer src.Close()

		imagePath := productName[0] + "_" + file.Filename + "_" + strconv.Itoa(int(timing))
		// log.Println("imagePath => ", imagePath)
		arrayImgPath = append(arrayImgPath, imagePath)
		clientStorage, err := storage.NewClient(ctx)
		if err != nil {
			log.Println("err storage.NewClient(ctx) => ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"Status": "internal error."})
		}

		bucket := clientStorage.Bucket(BUCKET)
		wc := bucket.Object(imagePath).NewWriter(ctx)
		_, err = io.Copy(wc, src)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"Status": "can't writer object"})
		}
		err = wc.Close()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"Status": "can't wc.Close"})
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
		ImagePath:    arrayImgPath,
	}

	// log.Println(payload)

	if _, err := clientDatastore.Put(ctx, key, &payload); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Status": "can't insert data"})
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

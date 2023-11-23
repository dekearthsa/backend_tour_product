package main

import (
	"log"
	"project_shopping_tour/service_product/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const port string = ":7700"

func main() {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE", "GET"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin,access-control-allow-headers"},
	}))

	router.GET("/api/debug", controller.ControllerDebug)
	router.POST("/api/sendfile", controller.ControllerInsertProduct)
	router.GET("/api/readfile", controller.ControllerReadProduct)
	router.POST("/api/update/product", controller.ControllerUpdateProduct)
	router.POST("/api/delete/product", controller.ControllerDeleteProduct)

	err := router.Run(port)
	if err != nil {
		log.Println("Service admin CRUD fail to start" + err.Error())
	}
	log.Println("Service admin CRUD start at port" + port + "debug => " + "http://localhost:" + port + "/debug")
}

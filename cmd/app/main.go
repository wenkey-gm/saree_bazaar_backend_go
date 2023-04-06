package main

import (
	"product_api/internal/adapters/handlers/sareehdl"
	"product_api/internal/adapters/repository/sareerepo"

	"product_api/internal/core/services"

	"github.com/gin-gonic/gin"
)

func main() {

	sareeRepository := sareerepo.NewSareeRepository()
	sareeService := services.NewSareeService(sareeRepository)

	sareeHandler := sareehdl.NewSareeHandler(sareeService)

	router := gin.New()
	router.GET("/sarees", sareeHandler.FindAll)
	router.GET("/sarees/:id", sareeHandler.Find)
	router.POST("/sarees", sareeHandler.Save)
	router.PUT("/sarees/:id", sareeHandler.Update)
	router.DELETE("/sarees/:id", sareeHandler.Delete)

	router.Run(":8080")

}

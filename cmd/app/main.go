package main

import (
	"product_api/internal/adapters/handlers/sareehdl"
	"product_api/internal/adapters/repository/saree_repo"
	"product_api/internal/utils"

	"product_api/internal/core/services"

	"github.com/gin-gonic/gin"
)

func main() {

	client := utils.DbConnection()

	// Saree Collection
	collection := utils.ConnectMongoDbCollection(client, utils.DB_NAME, utils.SAREE_COLLECTION)

	sareeRepository := saree_repo.NewSareeRepository(collection)
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

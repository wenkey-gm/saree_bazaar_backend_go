package main

import (
	"log"
	"product_api/internal/adapters/handlers/sareehdl"
	"product_api/internal/adapters/handlers/userhdl"
	"product_api/internal/adapters/repository/saree_repo"
	"product_api/internal/adapters/repository/user_repo"
	"product_api/internal/utils"

	"product_api/internal/core/services"

	"github.com/gin-gonic/gin"
)

func main() {

	var secretKey = []byte(utils.SECRET_KEY)

	client := utils.DbConnection()
	router := gin.New()

	// User Collection
	userCollection := utils.ConnectMongoDbCollection(client, utils.DB_NAME, utils.USER_COLLECTION)

	userRepository := user_repo.NewUserRepository(userCollection)
	userService := services.NewUserService(userRepository)

	userHandler := userhdl.NewUserHandler(userService)

	// Saree Collection
	sareeCollection := utils.ConnectMongoDbCollection(client, utils.DB_NAME, utils.SAREE_COLLECTION)

	sareeRepository := saree_repo.NewSareeRepository(sareeCollection)
	sareeService := services.NewSareeService(sareeRepository)

	sareeHandler := sareehdl.NewSareeHandler(sareeService)

	// Token Generator
	generator := utils.NewTokenGenerator(secretKey)
	token, _ := generator.GenerateToken()
	log.Println(token)

	router.GET("/users/:id", userHandler.Find)
	router.POST("/users", userHandler.Save)
	router.PUT("/users/:id", userHandler.Update)
	router.DELETE("/users/:id", userHandler.Delete)

	router.GET("/sarees", sareeHandler.FindAll)
	router.GET("/sarees/:id", sareeHandler.Find)
	router.POST("/sarees", sareeHandler.Save)
	router.PUT("/sarees/:id", sareeHandler.Update)
	router.DELETE("/sarees/:id", sareeHandler.Delete)

	router.Run(":8080")

}

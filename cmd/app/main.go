package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"os"
	"product_api/internal/adapters/handlers/sareehdl"
	"product_api/internal/adapters/handlers/userhdl"
	"product_api/internal/adapters/repository/saree_repo"
	"product_api/internal/adapters/repository/token_repo"
	"product_api/internal/adapters/repository/user_repo"
	"product_api/internal/utils"
	"strconv"

	"product_api/internal/core/services"

	"github.com/gin-gonic/gin"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	client := utils.DbConnection()
	router := gin.New()

	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := os.ReadFile(privKeyFile)
	if err != nil {
		fmt.Println("could not read private key file")
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := os.ReadFile(pubKeyFile)
	if err != nil {
		fmt.Println("could not read public key file")
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

	refreshSecret := os.Getenv("REFRESH_SECRET")

	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
	if err != nil {
		fmt.Errorf("could not parse ID_TOKEN_EXP as int: %w", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	tokenCollection := utils.ConnectMongoDbCollection(client, utils.DB_NAME, utils.TOKEN_COLLECTION)
	tokenRepository := token_repo.NewTokenRepository(tokenCollection)

	tokenService := services.NewTokenService(&services.TSConfig{
		TokenRepository:       tokenRepository,
		Pri:                   privKey,
		Pub:                   pubKey,
		RefreshSecret:         refreshSecret,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	})

	// User Collection
	userCollection := utils.ConnectMongoDbCollection(client, utils.DB_NAME, utils.USER_COLLECTION)

	userRepository := user_repo.NewUserRepository(userCollection)
	userService := services.NewUserService(userRepository)
	userHandler := userhdl.NewUserHandler(userService, tokenService)

	// Saree Collection
	sareeCollection := utils.ConnectMongoDbCollection(client, utils.DB_NAME, utils.SAREE_COLLECTION)

	sareeRepository := saree_repo.NewSareeRepository(sareeCollection)
	sareeService := services.NewSareeService(sareeRepository)
	sareeHandler := sareehdl.NewSareeHandler(sareeService)

	// Token Generator

	router.POST("/signup", userHandler.SignUp)
	router.POST("/login", userHandler.Login)

	router.GET("/sarees", sareeHandler.FindAll)
	router.GET("/sarees/:id", sareeHandler.Find)
	router.POST("/sarees", sareeHandler.Save)
	router.PUT("/sarees/:id", sareeHandler.Update)
	router.DELETE("/sarees/:id", sareeHandler.Delete)

	router.Run(":8080")

}

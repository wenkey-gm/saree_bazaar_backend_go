package main

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"product_api/internal/adapters/handlers/adminui"
	"product_api/internal/adapters/handlers/middlewares"
	"product_api/internal/adapters/handlers/sareehdl"
	"product_api/internal/adapters/handlers/uploadhdl"
	"product_api/internal/adapters/handlers/userhdl"
	"product_api/internal/adapters/repository/saree_repo"
	"product_api/internal/adapters/repository/token_repo"
	"product_api/internal/adapters/repository/user_repo"
	"product_api/internal/utils"
	"strconv"
	"strings"

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
	router.Use(middlewares.CORS(os.Getenv("ALLOWED_ORIGINS")))

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
		log.Printf("could not parse ID_TOKEN_EXP as int: %v", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		log.Printf("could not parse REFRESH_TOKEN_EXP as int: %v", err)
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
	userService := services.NewUserService(userRepository, strings.Split(os.Getenv("ADMIN_EMAILS"), ","))
	userHandler := userhdl.NewUserHandler(userService, tokenService)

	// Saree Collection
	sareeCollection := utils.ConnectMongoDbCollection(client, utils.DB_NAME, utils.SAREE_COLLECTION)

	sareeRepository := saree_repo.NewSareeRepository(sareeCollection)
	if n, err := sareeRepository.BackfillUIDs(context.Background()); err != nil {
		log.Printf("could not backfill saree IDs: %v", err)
	} else if n > 0 {
		log.Printf("assigned IDs to %d saree(s) saved without one", n)
	}
	sareeService := services.NewSareeService(sareeRepository)
	sareeHandler := sareehdl.NewSareeHandler(sareeService)

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	uploadHandler, err := uploadhdl.NewUploadHandler(uploadDir, os.Getenv("PUBLIC_BASE_URL"))
	if err != nil {
		log.Fatalf("could not create upload directory %q: %v", uploadDir, err)
	}

	requireAdmin := []gin.HandlerFunc{middlewares.AuthUser(tokenService), middlewares.RequireAdmin()}
	withAdmin := func(h gin.HandlerFunc) []gin.HandlerFunc {
		return append(append([]gin.HandlerFunc{}, requireAdmin...), h)
	}

	router.POST("/signup", userHandler.SignUp)
	router.POST("/login", userHandler.Login)
	router.DELETE("/signout", middlewares.AuthUser(tokenService), userHandler.SignOut)

	// The catalog is public so the storefront can list it; changes need an admin
	// (a user whose email is in ADMIN_EMAILS).
	router.GET("/sarees", sareeHandler.FindAll)
	router.GET("/sarees/:id", sareeHandler.Find)
	router.POST("/sarees", withAdmin(sareeHandler.Save)...)
	router.PUT("/sarees/:id", withAdmin(sareeHandler.Update)...)
	router.DELETE("/sarees/:id", withAdmin(sareeHandler.Delete)...)

	router.POST("/uploads", withAdmin(uploadHandler.Upload)...)
	router.GET("/uploads/:name", uploadHandler.Serve)

	// Catalog admin panel.
	adminui.Register(router)
	router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/admin") })

	router.Run(":8080")
}

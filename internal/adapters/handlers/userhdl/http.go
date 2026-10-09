package userhdl

import (
	"errors"
	"log"
	"net/http"
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService  ports.IUserService
	tokenService ports.ITokenService
}

func NewUserHandler(userService ports.IUserService, tokenService ports.ITokenService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

func userError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrEmailTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		log.Printf("user request failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
}

func (u *UserHandler) respondWithTokens(c *gin.Context, status int, user domain.User) {
	tokens, err := u.tokenService.GenerateTokens(c, &user, "")
	if err != nil {
		userError(c, err)
		return
	}
	c.JSON(status, gin.H{
		"tokens": tokens,
		"user":   gin.H{"id": user.ID, "email": user.Email, "role": user.Role},
	})
}

func (u *UserHandler) SignUp(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	created, err := u.userService.SignUp(user)
	if err != nil {
		userError(c, err)
		return
	}
	u.respondWithTokens(c, http.StatusCreated, created)
}

func (u *UserHandler) Login(c *gin.Context) {
	var user domain.SignRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	fetchedUser, err := u.userService.Login(user)
	if err != nil {
		userError(c, err)
		return
	}
	u.respondWithTokens(c, http.StatusOK, fetchedUser)
}

func (u *UserHandler) SignOut(c *gin.Context) {
	user := c.MustGet("user")

	ctx := c.Request.Context()
	if err := u.tokenService.SignOut(ctx, user.(*domain.User).ID.String()); err != nil {
		userError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user signed out successfully!",
	})
}

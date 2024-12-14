package userhdl

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"
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

func (u *UserHandler) SignUp(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, err)
		return
	}
	err := u.userService.SignUp(user)
	if err != nil {
		c.JSON(500, err)
		return
	}
	tokens, err := u.tokenService.GenerateTokens(c, &user, "")
	c.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}

func (u *UserHandler) Login(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, err)
		return
	}
	token, err := u.userService.Login(user)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, token)
}

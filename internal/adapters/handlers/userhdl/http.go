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
	var user domain.SignRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, err)
		return
	}
	fetchedUser, err := u.userService.Login(user)
	if err != nil {
		c.JSON(500, err)
		return
	}
	tokens, err := u.tokenService.GenerateTokens(c, &fetchedUser, "")

	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, gin.H{
		"tokens": tokens,
	})
}

func (u *UserHandler) SignOut(c *gin.Context) {
	user := c.MustGet("user")

	ctx := c.Request.Context()
	if err := u.tokenService.SignOut(ctx, user.(*domain.User).ID.String()); err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user signed out successfully!",
	})
}

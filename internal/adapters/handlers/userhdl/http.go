package userhdl

import (
	"github.com/gin-gonic/gin"
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"
)

type UserHandler struct {
	userService ports.IUserService
}

func NewUserHandler(userService ports.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) SignUp(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, err)
		return
	}
	savedUser, err := u.userService.SignUp(user)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, savedUser)
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

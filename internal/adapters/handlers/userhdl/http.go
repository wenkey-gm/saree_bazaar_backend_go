package userhdl

import (
	"github.com/gin-gonic/gin"
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) FindAll(c *gin.Context) {
	users, err := u.userService.FindAll()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, users)
}

func (u *UserHandler) Find(c *gin.Context) {
	id := c.Param("id")
	user, err := u.userService.Find(id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, user)
}

func (u *UserHandler) Save(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, err)
		return
	}
	savedUser, err := u.userService.Save(user)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, savedUser)
}

func (u *UserHandler) Update(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, err)
		return
	}
	id := c.Param("id")
	updatedUser, err := u.userService.Update(id, user)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, updatedUser)
}

func (u *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := u.userService.Delete(id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

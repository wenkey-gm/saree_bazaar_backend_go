package services

import (
	"product_api/internal/core/ports"
)

type UserService struct {
	repository ports.IUserRepository
}

func NewUserService(repository ports.IUserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

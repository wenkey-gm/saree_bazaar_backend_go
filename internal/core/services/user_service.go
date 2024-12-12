package services

import (
	"product_api/internal/core/domain"
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

func (s *UserService) SignUp(user domain.User) (domain.User, error) {
	return s.repository.Create(user)
}

func (s *UserService) Login(user domain.User) (domain.User, error) {
	return s.repository.Find(user.ID)
}

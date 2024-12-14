package services

import (
	"log"
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"
	"product_api/internal/utils"
)

type UserService struct {
	repository ports.IUserRepository
}

func NewUserService(repository ports.IUserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) SignUp(user domain.User) error {

	pw, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	user.Password = pw
	if err := s.repository.Create(user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(user domain.User) (domain.User, error) {
	return s.repository.Find(user.ID)
}

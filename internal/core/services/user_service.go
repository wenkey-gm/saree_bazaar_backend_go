package services

import (
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"
)

type UserService struct {
	repository ports.UserRepository
}

func NewUserService(repository ports.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (u *UserService) FindAll() ([]domain.User, error) {
	return u.repository.FindAll()
}

func (u *UserService) Find(id string) (domain.User, error) {
	return u.repository.Find(id)
}

func (u *UserService) Save(user domain.User) (domain.User, error) {
	return u.repository.Save(user)
}

func (u *UserService) Update(id string, user domain.User) (domain.User, error) {
	return u.repository.Update(id, user)
}

func (u *UserService) Delete(id string) error {
	return u.repository.Delete(id)
}

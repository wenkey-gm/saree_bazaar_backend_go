package ports

import (
	"github.com/google/uuid"
	"product_api/internal/core/domain"
)

type IUserRepository interface {
	Find(id uuid.UUID) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	Create(user domain.User) error
	Update(id string, user domain.User) (domain.User, error)
	Delete(id string) error
}

type IUserService interface {
	SignUp(user domain.User) error
	Login(user domain.SignRequest) (domain.User, error)
}

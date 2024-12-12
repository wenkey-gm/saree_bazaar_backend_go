package ports

import (
	"github.com/google/uuid"
	"product_api/internal/core/domain"
)

type IUserRepository interface {
	Find(id uuid.UUID) (domain.User, error)
	Create(user domain.User) (domain.User, error)
	Update(id string, user domain.User) (domain.User, error)
	Delete(id string) error
}

type IUserService interface {
	SignUp(user domain.User) (domain.User, error)
	Login(user domain.User) (domain.User, error)
}

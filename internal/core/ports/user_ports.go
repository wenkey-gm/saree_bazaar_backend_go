package ports

import "product_api/internal/core/domain"

type IUserRepository interface {
	Find(id string) (domain.User, error)
	Create(user domain.User) (domain.User, error)
	Update(id string, user domain.User) (domain.User, error)
	Delete(id string) error
}

type IUserService interface {
	Find(id string) (domain.User, error)
	Create(saree domain.User) (domain.User, error)
	Update(id string, saree domain.User) (domain.User, error)
	Delete(id string) error
}

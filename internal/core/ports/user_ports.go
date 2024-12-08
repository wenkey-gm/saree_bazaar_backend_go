package ports

import "product_api/internal/core/domain"

type UserRepository interface {
	FindAll() ([]domain.User, error)
	Find(id string) (domain.User, error)
	Save(user domain.User) (domain.User, error)
	Update(id string, user domain.User) (domain.User, error)
	Delete(id string) error
}

type UserService interface {
	FindAll() ([]domain.User, error)
	Find(id string) (domain.User, error)
	Save(saree domain.User) (domain.User, error)
	Update(id string, saree domain.User) (domain.User, error)
	Delete(id string) error
}

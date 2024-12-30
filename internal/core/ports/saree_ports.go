package ports

import "product_api/internal/core/domain"

type ISareeRepository interface {
	FindAll() ([]domain.Saree, error)
	Find(id string) (domain.Saree, error)
	Save(saree domain.Saree) (domain.Saree, error)
	Update(id string, saree domain.Saree) (domain.Saree, error)
	Delete(id string) error
}

type ISareeService interface {
	FindAll() ([]domain.Saree, error)
	Find(id string) (domain.Saree, error)
	Save(saree domain.Saree) (domain.Saree, error)
	Update(id string, saree domain.Saree) (domain.Saree, error)
	Delete(id string) error
}

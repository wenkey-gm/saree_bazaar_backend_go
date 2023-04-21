package ports

import "product_api/internal/core/domain"

type SareeRepository interface {
	FindAll() ([]domain.Saree, error)
	Find(id string) (domain.Saree, error)
	Save(saree domain.Saree) (domain.Saree, error)
	Update(id string, saree domain.Saree) (domain.Saree, error)
	Delete(id string) error
}

type SareeService interface {
	FindAll() ([]domain.Saree, error)
	Find(id string) (domain.Saree, error)
	Save(saree domain.Saree) (domain.Saree, error)
	Update(id string, saree domain.Saree) (domain.Saree, error)
	Delete(id string) error
}

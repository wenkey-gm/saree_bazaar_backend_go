package repository

import "saree_bazaar.com/pkg/domain/modal"

type SareeRepository interface {
	GetAllSarees() ([]modal.Saree, error)
	GetSaree(id string) (modal.Saree, error)
	CreateSaree(saree modal.Saree) error
	UpdateSaree(id string, saree modal.Saree) error
	DeleteSaree(id string) error
}

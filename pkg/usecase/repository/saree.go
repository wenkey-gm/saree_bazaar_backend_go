package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"saree_bazaar.com/pkg/domain/modal"
)

type SareeRepository interface {
	GetAllSarees() ([]modal.Saree, error)
	GetSaree(id primitive.ObjectID) (modal.Saree, error)
	CreateSaree(saree modal.Saree) (interface{}, error)
	UpdateSaree(id primitive.ObjectID, saree modal.Saree) error
	DeleteSaree(id primitive.ObjectID) (interface{}, error)
}

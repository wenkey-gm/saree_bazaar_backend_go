package services

import (
	"fmt"
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"
	"strings"

	"github.com/google/uuid"
)

type SareeService struct {
	repository ports.ISareeRepository
}

func NewSareeService(repository ports.ISareeRepository) *SareeService {
	return &SareeService{
		repository: repository,
	}
}

func validateSaree(saree *domain.Saree) error {
	saree.Name = strings.TrimSpace(saree.Name)
	if saree.Name == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if saree.Price < 0 {
		return fmt.Errorf("%w: price cannot be negative", domain.ErrInvalidInput)
	}
	if saree.Stock < 0 {
		return fmt.Errorf("%w: stock cannot be negative", domain.ErrInvalidInput)
	}
	if saree.ImagesUrl == nil {
		saree.ImagesUrl = []string{}
	}
	return nil
}

func (s *SareeService) FindAll() ([]domain.Saree, error) {
	return s.repository.FindAll()
}

func (s *SareeService) Find(id string) (domain.Saree, error) {
	return s.repository.Find(id)
}

func (s *SareeService) Save(saree domain.Saree) (domain.Saree, error) {
	if err := validateSaree(&saree); err != nil {
		return domain.Saree{}, err
	}
	saree.UID = uuid.New()
	return s.repository.Save(saree)
}

func (s *SareeService) Update(id string, saree domain.Saree) (domain.Saree, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.Saree{}, domain.ErrNotFound
	}
	if err := validateSaree(&saree); err != nil {
		return domain.Saree{}, err
	}
	// The ID comes from the URL; a body can't move a saree to another ID.
	saree.UID = uid
	return s.repository.Update(id, saree)
}

func (s *SareeService) Delete(id string) error {
	return s.repository.Delete(id)
}

package services

import (
	"product_api/internal/core/domain"
	"product_api/internal/core/ports"
)

type SareeService struct {
	repository ports.SareeRepository
}

func NewSareeService(repository ports.SareeRepository) *SareeService {
	return &SareeService{
		repository: repository,
	}
}

func (s *SareeService) FindAll() ([]domain.Saree, error) {
	return s.repository.FindAll()
}

func (s *SareeService) Find(id string) (domain.Saree, error) {
	return s.repository.Find(id)
}

func (s *SareeService) Save(saree domain.Saree) (domain.Saree, error) {
	return s.repository.Save(saree)
}

func (s *SareeService) Update(id string, saree domain.Saree) (domain.Saree, error) {
	return s.repository.Update(id, saree)
}

func (s *SareeService) Delete(id string) error {
	return s.repository.Delete(id)
}

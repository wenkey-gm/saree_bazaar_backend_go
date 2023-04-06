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

	sarees, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}
	return sarees, nil
}

func (s *SareeService) Find(id string) (domain.Saree, error) {
	saree, err := s.repository.Find(id)
	if err != nil {
		return domain.Saree{}, err
	}
	return saree, nil
}

func (s *SareeService) Save(saree domain.Saree) (domain.Saree, error) {
	saree, err := s.repository.Save(saree)
	if err != nil {
		return domain.Saree{}, err
	}
	return saree, nil
}

func (s *SareeService) Update(id string, saree domain.Saree) (domain.Saree, error) {
	saree, err := s.repository.Update(id, saree)
	if err != nil {
		return domain.Saree{}, err
	}
	return saree, nil
}

func (s *SareeService) Delete(id string) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

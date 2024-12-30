package saree_repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"product_api/internal/core/domain"
)

type SareeRepository struct {
	repo *mongo.Collection
}

func NewSareeRepository(collection *mongo.Collection) *SareeRepository {
	return &SareeRepository{
		repo: collection,
	}
}

func (s *SareeRepository) FindAll() ([]domain.Saree, error) {
	var sarees []domain.Saree
	cur, err := s.repo.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.Background()) {
		var saree domain.Saree
		err := cur.Decode(&saree)
		if err != nil {
			return nil, err
		}
		sarees = append(sarees, saree)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.Background())

	return sarees, nil
}

func (s *SareeRepository) Find(id string) (domain.Saree, error) {
	var saree domain.Saree
	objectId, objectErr := primitive.ObjectIDFromHex(id)
	if objectErr != nil {
		return domain.Saree{}, objectErr
	}
	err := s.repo.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&saree)
	if err != nil {
		return domain.Saree{}, err
	}
	return saree, nil
}

func (s *SareeRepository) Save(saree domain.Saree) (domain.Saree, error) {

	_, err := s.repo.InsertOne(context.Background(), saree)
	if err != nil {
		return domain.Saree{}, err
	}
	return saree, nil
}

func (s *SareeRepository) Update(id string, saree domain.Saree) (domain.Saree, error) {
	_, err := s.repo.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{"$set": saree})
	if err != nil {
		return domain.Saree{}, err
	}
	return saree, nil

}

func (s *SareeRepository) Delete(id string) error {
	objectId, objectErr := primitive.ObjectIDFromHex(id)
	if objectErr != nil {
		return objectErr
	}
	_, err := s.repo.DeleteOne(context.Background(), bson.M{"id": objectId})
	if err != nil {
		return err
	}
	return nil
}

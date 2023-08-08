package sareerepo

import (
	"context"
	"fmt"
	"log"
	"os"
	"product_api/internal/core/domain"
	"product_api/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDbCollection() *mongo.Collection {

	MONGO_URL := os.Getenv("MONGO_URL")
	if MONGO_URL == "" {
		MONGO_URL = utils.MONGO_URL
	}
	clientOptions := options.Client().ApplyURI(MONGO_URL) // mongodb://localhost:27017

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client.Database(utils.DB_NAME).Collection(utils.COLLECTION_NAME)
}

type sareeRepository struct {
	repo *mongo.Collection
}

func NewSareeRepository(collection *mongo.Collection) *sareeRepository {
	return &sareeRepository{
		repo: collection,
	}
}

func (s *sareeRepository) FindAll() ([]domain.Saree, error) {
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

func (s *sareeRepository) Find(id string) (domain.Saree, error) {
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

func (s *sareeRepository) Save(saree domain.Saree) (domain.Saree, error) {

	_, err := s.repo.InsertOne(context.Background(), saree)
	if err != nil {
		return domain.Saree{}, err
	}
	return saree, nil
}

func (s *sareeRepository) Update(id string, saree domain.Saree) (domain.Saree, error) {
	_, err := s.repo.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{"$set": saree})
	if err != nil {
		return domain.Saree{}, err
	}
	return saree, nil

}

func (s *sareeRepository) Delete(id string) error {
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

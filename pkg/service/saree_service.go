package service

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"saree_bazaar.com/pkg/domain/modal"
	"saree_bazaar.com/pkg/infrastructure/datastore"
	"saree_bazaar.com/pkg/usecase/repository"
)

type sareeService struct {
}

var db = datastore.ConnectDB()

// CreateSaree implements repository.SareeRepository
func (*sareeService) CreateSaree(saree modal.Saree) (interface{}, error) {

	response, err := db.Collection("sarees").InsertOne(context.TODO(), saree)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteSaree implements repository.SareeRepository
func (*sareeService) DeleteSaree(id primitive.ObjectID) (interface{}, error) {
	deleteResult, err := db.Collection("sarees").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return deleteResult, nil
}

// GetAllSarees implements repository.SareeRepository
func (*sareeService) GetAllSarees() ([]modal.Saree, error) {
	var sarees []modal.Saree

	cur, err := db.Collection("sarees").Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var saree modal.Saree
		err := cur.Decode(&saree)
		if err != nil {
			return nil, err
		}
		sarees = append(sarees, saree)
	}

	if err := cur.Err(); err != nil {

		return nil, err
	}

	return sarees, nil

}

// GetSaree implements repository.SareeRepository
func (*sareeService) GetSaree(id primitive.ObjectID) (modal.Saree, error) {
	var saree modal.Saree

	filter := bson.M{"_id": id}
	err := db.Collection("sarees").FindOne(context.TODO(), filter).Decode(&saree)

	if err != nil {
		return modal.Saree{}, err
	}

	return saree, nil
}

// UpdateSaree implements repository.SareeRepository
func (*sareeService) UpdateSaree(id primitive.ObjectID, saree modal.Saree) error {

	filter := bson.M{"_id": id}
	update := bson.D{
		{"$set", bson.D{
			{"name", saree.Name},
			{"price", saree.Price},
			{"image", saree.Image},
			{"type", saree.Type},
			{"color", saree.Color},
		}},
	}

	err := db.Collection("sarees").FindOneAndUpdate(context.TODO(), filter, update).Decode(&saree)

	if err != nil {
		return err
	}

	saree.ID = id

	return nil
}

// This function crates a new instance of sareeService which implements SareeRepository
func NewSareeService() repository.SareeRepository {
	return &sareeService{}
}

func GetError(err error) modal.ErrorResponse {
	response := modal.ErrorResponse{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: err.Error(),
	}

	return response
}

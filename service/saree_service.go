package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"saree_bazaar.com/config"
	"saree_bazaar.com/modal"
)

// GetAllSarees is a function
func GetAllSarees(w http.ResponseWriter, r *http.Request) {

	var sarees []modal.Saree

	collection := config.ConnectDB()

	cur, err := collection.Find(r.Context(), bson.M{})
	if err != nil {
		config.GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var saree modal.Saree
		err := cur.Decode(&saree)
		if err != nil {
			config.GetError(err, w)
			return
		}
		sarees = append(sarees, saree)
	}

	if err := cur.Err(); err != nil {
		config.GetError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sarees)
}

// GetSaree is a function

func GetSaree(w http.ResponseWriter, r *http.Request) {

	var saree modal.Saree

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}
	err := config.ConnectDB().FindOne(context.TODO(), filter).Decode(&saree)

	if err != nil {
		config.GetError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saree)
}

// CreateSaree is a function

func CreateSaree(w http.ResponseWriter, r *http.Request) {

	var saree modal.Saree

	_ = json.NewDecoder(r.Body).Decode(&saree)

	result, err := config.ConnectDB().InsertOne(context.TODO(), saree)

	if err != nil {
		config.GetError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// UpdateSaree is a function

func UpdateSaree(w http.ResponseWriter, r *http.Request) {

	var saree modal.Saree

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}

	_ = json.NewDecoder(r.Body).Decode(&saree)

	update := bson.D{
		{"$set", bson.D{
			{"name", saree.Name},
			{"price", saree.Price},
			{"image", saree.Image},
			{"type", saree.Type},
			{"color", saree.Color},
		}},
	}

	err := config.ConnectDB().FindOneAndUpdate(context.TODO(), filter, update).Decode(&saree)

	if err != nil {
		config.GetError(err, w)
		return
	}

	saree.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saree)
}

// DeleteSaree is a function

func DeleteSaree(w http.ResponseWriter, r *http.Request) {

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	deleteResult, err := config.ConnectDB().DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		config.GetError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deleteResult)
}

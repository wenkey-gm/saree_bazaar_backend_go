package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"saree_bazaar.com/pkg/domain/modal"
	"saree_bazaar.com/pkg/infrastructure/datastore"
	"saree_bazaar.com/pkg/service"
)

var db = datastore.ConnectDB()

func GetAllSarees(w http.ResponseWriter, r *http.Request) {

	response, err := service.NewSareeService().GetAllSarees()
	if err != nil {
		GetError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSaree is a function

func GetSaree(w http.ResponseWriter, r *http.Request) {

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	response, err := service.NewSareeService().GetSaree(id)
	if err != nil {
		GetError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateSaree is a function

func CreateSaree(w http.ResponseWriter, r *http.Request) {

	var saree modal.Saree

	_ = json.NewDecoder(r.Body).Decode(&saree)

	result, err := service.NewSareeService().CreateSaree(saree)
	if err != nil {
		GetError(err, w)
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

	_ = json.NewDecoder(r.Body).Decode(&saree)

	err := service.NewSareeService().UpdateSaree(id, saree)
	if err != nil {
		GetError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Saree updated successfully")
}

// DeleteSaree is a function

func DeleteSaree(w http.ResponseWriter, r *http.Request) {

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	deleteResult, err := service.NewSareeService().DeleteSaree(id)
	if err != nil {
		GetError(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deleteResult)
}

func GetError(err error, w http.ResponseWriter) {
	response := service.GetError(err)
	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

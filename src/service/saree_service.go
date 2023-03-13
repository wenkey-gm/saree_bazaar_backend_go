package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"saree_bazaar.com/src/config"
	"saree_bazaar.com/src/modal"
)

// GetAllSarees is a function
func GetAllSarees(w http.ResponseWriter, r *http.Request) {
	sarees, text := config.ConnectDB().Database("saree").Collection("sarees").Find(r.Context(), modal.Saree{})

	fmt.Println(text)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sarees)
}

// GetSaree is a function

func GetSaree(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	saree := modal.Saree{ID: id, Name: "Saree " + strconv.Itoa(id), Price: id * 100}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saree)
}

// CreateSaree is a function

func CreateSaree(w http.ResponseWriter, r *http.Request) {

	var saree modal.Saree
	json.NewDecoder(r.Body).Decode(&saree)

	fmt.Println(saree)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saree)
}

// UpdateSaree is a function

func UpdateSaree(w http.ResponseWriter, r *http.Request) {

	var saree modal.Saree
	json.NewDecoder(r.Body).Decode(&saree)

	fmt.Println(saree)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saree)
}

// DeleteSaree is a function

func DeleteSaree(w http.ResponseWriter, r *http.Request) {

	var saree modal.Saree
	json.NewDecoder(r.Body).Decode(&saree)

	fmt.Println(saree)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saree)
}

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	controller "saree_bazaar.com/pkg/controller"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/sarees", controller.GetAllSarees).Methods("GET")

	r.HandleFunc("/saree/{id}/", controller.GetSaree).Methods("GET")

	r.HandleFunc("/create", controller.CreateSaree).Methods("POST")

	r.HandleFunc("/update/{id}/", controller.UpdateSaree).Methods("PUT")

	r.HandleFunc("/delete/{id}/", controller.DeleteSaree).Methods("DELETE")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}

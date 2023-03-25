package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"saree_bazaar.com/pkg/config"
	controller "saree_bazaar.com/pkg/controller"
)

func main() {

	config.ReadConfig()

	r := mux.NewRouter()

	r.HandleFunc("/", controller.GetAllSarees)

	r.HandleFunc("/saree", controller.GetSaree)

	r.HandleFunc("/create", controller.CreateSaree)

	r.HandleFunc("/update", controller.UpdateSaree)

	r.HandleFunc("/delete", controller.DeleteSaree)

	log.Fatal(http.ListenAndServe(":8000", r))
}

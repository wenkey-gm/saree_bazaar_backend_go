package main

import (
	"net/http"

	controller "saree_bazaar.com/src/service"
)

func main() {

	http.HandleFunc("/", controller.GetAllSarees)

	http.HandleFunc("/saree", controller.GetSaree)

	http.HandleFunc("/create", controller.CreateSaree)

	http.HandleFunc("/update", controller.UpdateSaree)

	http.HandleFunc("/delete", controller.DeleteSaree)

	http.ListenAndServe(":8080", nil)
}

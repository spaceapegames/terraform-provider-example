package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"terraform-provider-blog/item"
)

func main() {
	// TODO Load from file
	itemService := item.NewItemService()

	r := mux.NewRouter()
	r.HandleFunc("/item", itemService.GetItems).Methods("GET")
	r.HandleFunc("/item", itemService.PostItem).Methods("POST")
	r.HandleFunc("/item/{id}", itemService.GetItem).Methods("GET")
	r.HandleFunc("/item/{id}", itemService.PutItem).Methods("PUT")
	//r.HandleFunc("/item/{id}", HomeHandler).Methods("DELETE")

	log.Println("Starting server on localhost:3001")
	err := http.ListenAndServe("localhost:3001", r)
	if err != nil {
		log.Fatal(err)
	}
}
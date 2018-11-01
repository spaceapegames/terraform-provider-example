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
	r.HandleFunc("/item", auth(itemService.GetItems)).Methods("GET")
	r.HandleFunc("/item", auth(itemService.PostItem)).Methods("POST")
	r.HandleFunc("/item/{id}", auth(itemService.GetItem)).Methods("GET")
	r.HandleFunc("/item/{id}", auth(itemService.PutItem)).Methods("PUT")
	r.HandleFunc("/item/{id}", auth(itemService.DeleteItem)).Methods("DELETE")

	log.Println("Starting server on localhost:3001")
	err := http.ListenAndServe("localhost:3001", r)
	if err != nil {
		log.Fatal(err)
	}
}

func auth(handlerFunc http.HandlerFunc) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Please supply and Authorization token", http.StatusUnauthorized)
			return
		}
		handlerFunc(w, r)
		return
	}
}
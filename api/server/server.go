package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

type Service struct {
	connectionString string
	items            map[string]Item
	sync.RWMutex
}

func NewService(connectionString string, items map[string]Item) *Service {
	return &Service{
		connectionString: connectionString,
		items: items,
	}
}


func (s *Service) ListenAndServe() error {
	r := mux.NewRouter()
	r.HandleFunc("/item", auth(s.PostItem)).Methods("POST")
	r.HandleFunc("/item", auth(s.GetItems)).Methods("GET")
	r.HandleFunc("/item/{name}", auth(s.GetItem)).Methods("GET")
	r.HandleFunc("/item/{name}", auth(s.PutItem)).Methods("PUT")
	r.HandleFunc("/item/{name}", auth(s.DeleteItem)).Methods("DELETE")

	log.Printf("Starting server on %s", s.connectionString)
	err := http.ListenAndServe(s.connectionString, r)
	if err != nil {
		return err
	}
	return nil
}

// auth checks that a non-empty authorization header has been sent with the request
func auth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method)
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Please supply and Authorization token", http.StatusUnauthorized)
			return
		}
		handlerFunc(w, r)
		return
	}
}

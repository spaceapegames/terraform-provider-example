package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"regexp"
)

// Item represents a single Item
type Item struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// GetItems returns all of the Items that exist in the server
func (s *Service) GetItems(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	defer s.RUnlock()
	s.shuffleItemTags()
	err := json.NewEncoder(w).Encode(s.items)
	if err != nil {
		log.Println(err)
	}
}

// PostItem handles adding a new Item
func (s *Service) PostItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(item.Name)) {
		http.Error(w, "item names cannot contain whitespace", 400)
		return
	}

	s.Lock()
	defer s.Unlock()

	if s.itemExists(item.Name) {
		http.Error(w, fmt.Sprintf("item %s already exists", item.Name), http.StatusBadRequest)
		return
	}

	s.items[item.Name] = item
	log.Printf("added item: %s", item.Name)
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		log.Printf("error sending response - %s", err)
	}
}

// PutItem handles updating an Item with a specific name
func (s *Service) PutItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemName := vars["name"]
	if itemName == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var item Item
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	s.Lock()
	defer s.Unlock()

	if !s.itemExists(itemName) {
		log.Printf("item %s does not exist", itemName)
		http.Error(w, fmt.Sprintf("item %v does not exist", itemName), http.StatusBadRequest)
		return
	}

	s.items[itemName] = item
	log.Printf("updated item: %s", item.Name)
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		log.Printf("error sending response - %s", err)
	}
}

// DeleteItem handles removing an Item with a specific name
func (s *Service) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemName := vars["name"]
	if itemName == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	s.Lock()
	defer s.Unlock()

	if !s.itemExists(itemName) {
		http.Error(w, fmt.Sprintf("item %s does not exists", itemName), http.StatusNotFound)
		return
	}

	delete(s.items, itemName)

	_, err := fmt.Fprintf(w, "Deleted item with name %s", itemName)
	if err != nil {
		log.Println(err)
	}
}

// GetItem handles retrieving an Item with a specific name
func (s *Service) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemName := vars["name"]
	if itemName == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	s.RLock()
	defer s.RUnlock()
	s.shuffleItemTags()
	if !s.itemExists(itemName) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(s.items[itemName])
	if err != nil {
		log.Println(err)
		return
	}
}

// itemExists checks if an item exists in or not. Does not lock access to the itemService, expects this to
// be done by the calling method
func (s *Service) itemExists(itemName string) bool {
	if _, ok := s.items[itemName]; ok {
		return true
	}
	return false
}

// suffleItemTags shuffles the order of the tags within each item in the itemService.Does not lock access
// to the itemService, expects this to be done by the calling method
func (s *Service) shuffleItemTags() {
	for _, item := range s.items {
		for i := range item.Tags {
			j := rand.Intn(i + 1)
			item.Tags[i], item.Tags[j] = item.Tags[j], item.Tags[i]
		}
	}
}

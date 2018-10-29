package item

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        []Tag  `json:"tags"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type service struct {
	items map[int]Item
	count int
	sync.RWMutex
}

func NewItemService() *service {
	return &service{
		items: map[int]Item{},
	}
}

func NewItemServicePrePopulated(items map[int]Item) *service {
	return &service{
		items: items,
	}
}

func (s *service) GetItems(w http.ResponseWriter, r *http.Request) {
	log.Println("GetItems")
	s.Lock()
	defer s.Unlock()
	s.shuffleItemTags()
	err := json.NewEncoder(w).Encode(s.items)
	if err != nil {
		log.Println(err)
	}
}

func (s *service) PostItem(w http.ResponseWriter, r *http.Request) {
	log.Println("PostItems")
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
	for _, i := range s.items {
		log.Println(i.Name, item.Name)
		if i.Name == item.Name {
			http.Error(w, "item already exists", http.StatusBadRequest)
			return
		}
	}
	s.Lock()
	defer s.Unlock()
	s.count++
	s.items[s.count] = item
	log.Printf("added item: %s", item.Name)
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		log.Printf("error sending response - %s", err)
	}
}

func (s *service) PutItem(w http.ResponseWriter, r *http.Request) {
	log.Println("PutItems")

	vars := mux.Vars(r)
	itemId := vars["id"]
	if itemId == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	itemIdNum, err := strconv.Atoi(itemId)
	if err != nil {
		http.Error(w, fmt.Sprintf("expected and integer as an ID, got %s", itemId), http.StatusBadRequest)
		return
	}
	liveItem, ok := s.items[itemIdNum]
	if !ok {
		log.Printf("item %v does not exist", itemIdNum)
		http.Error(w, fmt.Sprintf("item %v does not exist", itemIdNum), http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&liveItem)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid json provided", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(&liveItem)
	if err != nil {
		log.Println(err)
	}
}

func (s *service) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemId := vars["id"]
	if itemId == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	s.Lock()
	defer s.Unlock()
	s.shuffleItemTags()
	itemIdint, err := strconv.Atoi(itemId)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	item := Item{}
	item, ok := s.items[itemIdint]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *service) shuffleItemTags() {
	for _, item := range s.items {
		for i := range item.Tags {
			j := rand.Intn(i + 1)
			item.Tags[i], item.Tags[j] = item.Tags[j], item.Tags[i]
		}
	}
}

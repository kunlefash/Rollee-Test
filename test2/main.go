package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type WordService struct {
	words map[string]int
	mutex sync.RWMutex
}

func NewWordService() *WordService {
	return &WordService{
		words: make(map[string]int),
	}
}

func (ws *WordService) AddWord(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Word string `json:"word"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	word := strings.TrimSpace(req.Word)
	if !isWordValid(word) {
		http.Error(w, "Invalid word format.", http.StatusBadRequest)
		return
	}

	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	ws.words[strings.ToLower(word)]++
}

func (ws *WordService) GetMostFrequentWord(w http.ResponseWriter, r *http.Request) {
	prefix := strings.ToLower(r.FormValue("prefix"))

	ws.mutex.RLock()
	words := ws.words
	ws.mutex.RUnlock()

	var mostFrequentWord string
	maxCount := 0

	for word, count := range words {
		if strings.HasPrefix(word, prefix) && count > maxCount {
			mostFrequentWord = word
			maxCount = count
		}
	}

	if mostFrequentWord == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"word": mostFrequentWord})
}

func isWordValid(word string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z]+$`, word)
	return match
}

func main() {
	wordService := NewWordService()

	router := mux.NewRouter()
	router.HandleFunc("/service/word", wordService.AddWord).Methods("POST")
	router.HandleFunc("/service/prefix", wordService.GetMostFrequentWord).Methods("GET")

	log.Println("The server is on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

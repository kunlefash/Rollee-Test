package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	storageService = &StorageService{
		store: make(map[string]int),
	}
)

// StorageService function represents the storage for words and their frequency.
type StorageService struct {
	sync.Mutex
	store map[string]int
}

func (s *StorageService) AddWord(word string) {
	s.Lock()
	defer s.Unlock()
	s.store[word]++
}

// GetMostFrequentWord function returns the most frequent word starting with the given prefix.
func (s *StorageService) GetMostFrequentWord(prefix string) (string, error) {
	s.Lock()
	defer s.Unlock()
	maxFreq := 0
	mostFrequentWord := ""
	for word, freq := range s.store {
		if strings.HasPrefix(word, prefix) && freq > maxFreq {
			maxFreq = freq
			mostFrequentWord = word
		}
	}
	if mostFrequentWord == "" {
		return "", errors.New("no word found with the given prefix")
	}
	return mostFrequentWord, nil
}

func main() {
	http.HandleFunc("/service/word", handleWord)
	http.HandleFunc("/service/prefix", handlePrefix)

	fmt.Println("HTTP Service started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	word := r.FormValue("word")
	isValid, err := regexp.MatchString(`^[a-zA-Z]+$`, word)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	if !isValid {
		http.Error(w, "Invalid word format. Only alphabetic characters are allowed", http.StatusNotAcceptable)
		return
	}

	storageService.AddWord(strings.ToLower(word))
	fmt.Fprintf(w, "Word '%s' stored\n", word)
}

func handlePrefix(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prefix := r.FormValue("prefix")
	isValid, err := regexp.MatchString(`^[a-zA-Z]+$`, prefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	if !isValid {
		http.Error(w, "Invalid prefix format. Only alphabetic characters are allowed", http.StatusNotAcceptable)
		return
	}

	mostFrequentWord, err := storageService.GetMostFrequentWord(strings.ToLower(prefix))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Most frequent word with prefix '%s': %s\n", prefix, mostFrequentWord)
}

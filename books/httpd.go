package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	maxBookSize = 1 << 20 //1MB
)

func main() {
	// GO ROUTER
	// built in router cand do two kind of matches:
	// exact match: when the route does not end with / "/health"
	// prefix match: when the route ends with / "/health/"
	// http.HandleFunc("/health", healthHandler)
	// http.HandleFunc("/add", addHandler)

	// GORILLLA ROUTER
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/books", addHandler).Methods("POST")
	r.HandleFunc("/books/{isbn}", getHandler).Methods("GET")
	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error: %s", err)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	// step 1. unmarshal & validate
	vars := mux.Vars(r)
	isbn := vars["isbn"] //should match the name on line 29 (getHandler router)
	if isbn == "" {
		http.Error(w, "missing ISBN", http.StatusBadRequest)
		return
	}

	// step2: work
	book, ok := getBook(isbn)
	if !ok {
		msg := fmt.Sprintf("unkown ISBN: %s", isbn)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// step 3: marshal response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Printf("json serialization failed: %s", err)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	//Step 1: Unmarshal & Validate
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	dec := json.NewDecoder(io.LimitReader(r.Body, maxBookSize))

	var book Book
	if err := dec.Decode(&book); err != nil {
		// TODO: Check we're not leaking sensitive data in err
		// vs the other times we send an error we are sending a string we write
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("here", book)
	// some validating of data -- framework for validation: cuelang
	if book.Author == "" || book.ISBN == "" || book.Title == "" {
		http.Error(w, "missing data", http.StatusBadRequest)
		return
	}

	// Step 2: Do the work (addBook cmes from db.go)
	count := addBook(book)

	// Step 3: marshal Reply
	reply := map[string]interface{}{
		"isbn":  book.ISBN,
		"count": count,
	}
	// TODO: Use json.Marshal & then write
	if err := json.NewEncoder(w).Encode(reply); err != nil {
		log.Printf("json serialization failed: %s", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK\n")
}

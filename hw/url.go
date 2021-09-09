package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/gorilla/mux"
)

var (
	urlDB   sync.Map
	counter uint64
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/url", handleURL).Methods("POST")
	r.HandleFunc("/s/{id}", getURL).Methods("GET")
	//http.Handle("/", r)

	addy := ":3030"
	log.Printf("server ready on %s", addy)
	if err := http.ListenAndServe(addy, r); err != nil {
		log.Fatalf("error: %s", err)
	}

}

func getURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	i, ok := urlDB.Load(id)
	if !ok {
		http.Error(w, fmt.Sprintf("%s not found", id), http.StatusNotFound)
		return
	}

	url, ok := i.(string)
	if !ok {
		urlDB.Delete(id)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func handleURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string
		ID  string
	}

	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/// ADDED THIS CHECK TO CONFIRM THE URL WAS CORRECT FORMAT
	if _, err := url.Parse(req.URL); err != nil {
		http.Error(w, "malformed URL", http.StatusBadRequest)
		return
	}

	i := atomic.AddUint64(&counter, 1)
	req.ID = base62Encode(i)
	urlDB.Store(req.ID, req.URL)

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(req); err != nil {
		log.Printf("serialization failed: %s", err)
	}

}

package main

import (
	"sync"
)

var (
	db   = make(map[string]Book) // ISBN -> Book
	lock sync.RWMutex
)

type Book struct {
	ISBN   string
	Author string
	Title  string
}

func addBook(book Book) int {
	lock.Lock()
	defer lock.Unlock()

	db[book.ISBN] = book
	return len(db)
}

func getBook(isbn string) (Book, bool) {
	lock.RLock()
	defer lock.RUnlock()

	book, ok := db[isbn]
	return book, ok
}

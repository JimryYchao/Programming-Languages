package main

import (
	"log"
	"net/http"
	"testing"
)

// ? go test -v -run=^TestSearch$
func TestSearch(t *testing.T) {
	http.HandleFunc("/search", handleSearch)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

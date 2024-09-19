package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewInMemoryPlayerStore()
	handler := NewPlayerServer(store)
	log.Fatal(http.ListenAndServe(":3000", handler))
}

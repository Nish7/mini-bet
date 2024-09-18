package main

import (
	"log"
	"net/http"
)

func main() {
	handler := &PlayerServer{NewInMemoryPlayerStore()}
	log.Fatal(http.ListenAndServe(":3000", handler))
}

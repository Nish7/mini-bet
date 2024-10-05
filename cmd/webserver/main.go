package main

import (
	"github.com/nish7/mini-bet"
	"log"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	handler := poker.NewPlayerServer(store)
	log.Fatal(http.ListenAndServe(":3000", handler))
}

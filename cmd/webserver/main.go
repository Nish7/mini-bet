package main

import (
	"github.com/nish7/mini-bet"
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, _ := poker.NewFileSystemPlayerStore(db)
	handler := poker.NewPlayerServer(store)
	log.Fatal(http.ListenAndServe(":3000", handler))
}

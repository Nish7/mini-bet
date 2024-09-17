package main

import (
	"fmt"
	"net/http"
	"strings"
)

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	// case insensitive path handling
	player := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/players/"))

	if player == "Pepper" {
		fmt.Fprint(w, "20")
		return
	}

	if player == "floyd" {
		fmt.Fprint(w, "10")
		return
	}
}

package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(database *os.File) *FileSystemPlayerStore {
	database.Seek(0, io.SeekStart)
	league, err := NewLeague(database)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	return &FileSystemPlayerStore{
		league:   league,
		database: json.NewEncoder(&Tape{database}),
	}
}

func (f *FileSystemPlayerStore) GetLeague() League {
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.GetLeague().Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWins(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.database.Encode(f.league)
}

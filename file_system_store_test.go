package main

import (
	"io"
	"os"
	"testing"
)

func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFunc := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFunc
}

func TestFileSystemStoreTest(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store := FileSystemPlayerStore{database}
		defer cleanDatabase()

		got := store.GetLeague()

		want := []Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assertLeague(t, got, want)

	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store := FileSystemPlayerStore{database}
		defer cleanDatabase()

		got := store.GetPlayerScore("Chris")

		want := 33

		if got != want {
			t.Errorf("got %d wanted %d", got, want)
		}

	})
}

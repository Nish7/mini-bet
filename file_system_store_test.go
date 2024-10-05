package main

import (
	"os"
	"testing"
)

func TestFileSystemStoreTest(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store, err := NewFileSystemPlayerStore(database)
		defer cleanDatabase()

		got := store.GetLeague()

		want := League{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)
		assertNoError(t, err)

	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store, err := NewFileSystemPlayerStore(database)
		defer cleanDatabase()

		got := store.GetPlayerScore("Chris")

		want := 33

		assertNoError(t, err)

		if got != want {
			t.Errorf("got %d wanted %d", got, want)
		}
	})

	t.Run("store wins for exisiting players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store, err := NewFileSystemPlayerStore(database)
		defer cleanDatabase()

		store.RecordWins("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34

		assertScoreEqual(t, got, want)
		assertNoError(t, err)
	})

	t.Run("store wins for exisiting players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store, _ := NewFileSystemPlayerStore(database)
		defer cleanDatabase()

		store.RecordWins("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1

		assertScoreEqual(t, got, want)
	})

	t.Run("works with the empty string", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()
		_, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()

		want := League{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
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

func assertScoreEqual(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("got %d and wanted %d", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

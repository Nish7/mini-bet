package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWins(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGetPlayers(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{
			"pepper": 20,
			"floyd":  10,
		},
	}

	server := &PlayerServer{store}

	t.Run("returns peppers score", func(t *testing.T) {
		request, _ := newGetScoreRequest("pepper")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertStatusCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)
	})

	t.Run("returns floyds score", func(t *testing.T) {
		request, _ := newGetScoreRequest("floyd")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"

		assertStatusCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request, _ := newGetScoreRequest("Appollo")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		assertStatusCode(t, got, want)
	})
}

func TestStoreScores(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{},
		nil,
	}

	server := &PlayerServer{store}

	t.Run("it returns accepted on POST", func(t *testing.T) {
		name := "Pepper"
		request := newPostWinRequest(name)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWins wanted %d calls", len(store.winCalls), 1)
		}

		if store.winCalls[0] != name {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], name)
		}

		assertStatusCode(t, response.Code, http.StatusAccepted)
	})
}

func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Got status code %d wanted %d", got, want)
	}
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetScoreRequest(name string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, "/players/"+name, nil)
	return req, err
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

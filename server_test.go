package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
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
		request, _ := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertResponseBody(t, got, want)
	})

	t.Run("returns floyds score", func(t *testing.T) {
		request, _ := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"

		assertResponseBody(t, got, want)
	})
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

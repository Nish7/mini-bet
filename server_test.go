package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWins(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func TestGetPlayers(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{
			"pepper": 20,
			"floyd":  10,
		},
	}

	server := NewPlayerServer(store)

	t.Run("returns peppers score", func(t *testing.T) {
		request := newGetScoreRequest("pepper")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertStatusCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)
	})

	t.Run("returns floyds score", func(t *testing.T) {
		request := newGetScoreRequest("floyd")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"

		assertStatusCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Appollo")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		assertStatusCode(t, got, want)
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the legaue table as JSON", func(t *testing.T) {
		league := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := &StubPlayerStore{league: league}
		server := NewPlayerServer(store)

		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, "/league", nil)

		server.ServeHTTP(response, request)

		var got League

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player , %v", response.Body, err)
		}

		if response.Result().Header.Get("content-type") != "application/json" {
			t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
		}

		assertStatusCode(t, response.Code, http.StatusOK)
		assertLeague(t, got, league)
	})
}

func TestStoreScores(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}

	server := NewPlayerServer(store)

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

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/players/"+name, nil)
	return req
}

func newGetLeague() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of player %v", body, err)
	}

	return
}

func assertLeague(t *testing.T, got, wantLeague []Player) {
	t.Helper()

	if !reflect.DeepEqual(got, wantLeague) {
		t.Errorf("got %v want %v", got, wantLeague)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

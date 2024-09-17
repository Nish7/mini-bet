package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPlayers(t *testing.T) {
	t.Run("returns peppers score", func(t *testing.T) {
		request, _ := newGetScoreRequest("pepper")
		response := httptest.NewRecorder()
		PlayerServer(response, request)

		got := response.Body.String()
		want := "20"

		assertResponseBody(t, got, want)
	})

	t.Run("returns floyds score", func(t *testing.T) {
		request, _ := newGetScoreRequest("floyd")
		response := httptest.NewRecorder()
		PlayerServer(response, request)

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

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordWinsAndRetrieveThem(t *testing.T) {
	server := PlayerServer{NewInMemoryPlayerStore()}
	player := "pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertStatusCode(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), "3")
}

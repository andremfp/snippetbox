package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	t.Run("/ gets 200 and hello message", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		homeHandler(response, request)

		want := "Hello from Snippetbox"

		assertResponseBody(t, response.Body.String(), want)
		assertResponseCode(t, response.Code, http.StatusOK)

	})

	t.Run("display snippet with id 1", func(t *testing.T) {
		id := 1
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/snippet/view?id=%d", id), nil)
		response := httptest.NewRecorder()

		snippetViewHandler(response, request)

		want := fmt.Sprintf("Display a specific snippet with ID %d...", id)

		assertResponseBody(t, response.Body.String(), want)
		assertResponseCode(t, response.Code, http.StatusOK)

	})

	t.Run("display snippet with invalid id gets 404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/snippet/view?id=abcdef", nil)
		response := httptest.NewRecorder()

		snippetViewHandler(response, request)

		want := "404 page not found\n"

		assertResponseBody(t, response.Body.String(), want)
		assertResponseCode(t, response.Code, http.StatusNotFound)

	})

	t.Run("/snippet/create POST gets 200 and hello message", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/snippet/create", nil)
		response := httptest.NewRecorder()

		snippetCreateHandler(response, request)

		want := "Create a new snippet..."

		assertResponseBody(t, response.Body.String(), want)
		assertResponseCode(t, response.Code, http.StatusOK)

	})

	t.Run("/snippet/create without POST gets a 405", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/snippet/create", nil)
		response := httptest.NewRecorder()

		snippetCreateHandler(response, request)

		want := "Method Not Allowed\n"

		gotAllowHeader := response.Header().Get("Allow")
		wantAllowHeader := "POST"

		if gotAllowHeader != wantAllowHeader {
			t.Errorf("got 'Allow' header %q, want %q", gotAllowHeader, wantAllowHeader)
		}

		assertResponseBody(t, response.Body.String(), want)
		assertResponseCode(t, response.Code, http.StatusMethodNotAllowed)

	})

	t.Run("anything else return 404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/abcdef", nil)
		response := httptest.NewRecorder()

		homeHandler(response, request)

		want := "404 page not found\n"

		assertResponseBody(t, response.Body.String(), want)
		assertResponseCode(t, response.Code, http.StatusNotFound)

	})
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got response %q, want %q", got, want)

	}
}

func assertResponseCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got response code %d, want %d", got, want)

	}
}

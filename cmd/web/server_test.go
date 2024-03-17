package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubSnippetStore struct{}

func (s *StubSnippetStore) Insert(title, content string, expires int) (int, error) {
	return 1, nil
}

func TestServer(t *testing.T) {

	testApp := &application{}
	testApp.snippetStore = &StubSnippetStore{}
	testServer := httptest.NewServer(testApp.NewServeMux())
	testClient := testServer.Client()
	testClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	defer testServer.Close()

	t.Run("root path returns 200", func(t *testing.T) {
		response, err := testClient.Get(fmt.Sprintf("%s/", testServer.URL))
		if err != nil {
			t.Fatalf("could not make request to test server, %v", err)
		}

		assertResponseCode(t, response.StatusCode, http.StatusOK)

	})

	t.Run("display snippet with id 1", func(t *testing.T) {
		id := 1
		response, err := testClient.Get(fmt.Sprintf("%s/snippet/view?id=%d", testServer.URL, id))
		if err != nil {
			t.Fatalf("could not make request to test server, %v", err)
		}
		defer response.Body.Close()

		got, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("could not read response body, %v", err)
		}

		want := fmt.Sprintf("Display a specific snippet with ID %d...", id)

		assertResponseBody(t, string(got), want)
		assertResponseCode(t, response.StatusCode, http.StatusOK)

	})

	t.Run("display snippet with invalid id returns 404", func(t *testing.T) {
		response, err := testClient.Get(fmt.Sprintf("%s/snippet/view?id=abcdef", testServer.URL))
		if err != nil {
			t.Fatalf("could not make request to test server, %v", err)
		}
		defer response.Body.Close()

		got, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("could not read response body, %v", err)
		}

		want := "Not Found\n"

		assertResponseBody(t, string(got), want)
		assertResponseCode(t, response.StatusCode, http.StatusNotFound)

	})

	t.Run("/snippet/create POST returns 303 and redirects to snippet view", func(t *testing.T) {
		response, err := testClient.Post(fmt.Sprintf("%s/snippet/create", testServer.URL), "", nil)
		if err != nil {
			t.Fatalf("could not make request to test server, %v", err)
		}
		defer response.Body.Close()

		gotRedirect := response.Header.Get("Location")
		wantRedirect := "/snippet/view?id=1"

		if gotRedirect != wantRedirect {
			t.Errorf("got redirect %s, want %s", gotRedirect, wantRedirect)
		}

		assertResponseCode(t, response.StatusCode, http.StatusSeeOther)

	})

	t.Run("/snippet/create without POST returns a 405", func(t *testing.T) {
		response, err := testClient.Get(fmt.Sprintf("%s/snippet/create", testServer.URL))
		if err != nil {
			t.Fatalf("could not make request to test server, %v", err)
		}
		defer response.Body.Close()

		got, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("could not read response body, %v", err)
		}

		want := "Method Not Allowed\n"

		gotAllowHeader := response.Header.Get("Allow")
		wantAllowHeader := "POST"

		if gotAllowHeader != wantAllowHeader {
			t.Errorf("got 'Allow' header %q, want %q", gotAllowHeader, wantAllowHeader)
		}

		assertResponseBody(t, string(got), want)
		assertResponseCode(t, response.StatusCode, http.StatusMethodNotAllowed)

	})

	t.Run("/static/ returns 200", func(t *testing.T) {
		response, err := testClient.Get(fmt.Sprintf("%s/static/", testServer.URL))
		if err != nil {
			t.Fatalf("could not make request to test server, %v", err)
		}
		defer response.Body.Close()

		assertResponseCode(t, response.StatusCode, http.StatusOK)

	})

	t.Run("anything else return 404", func(t *testing.T) {
		response, err := testClient.Get(fmt.Sprintf("%s/abcdef", testServer.URL))
		if err != nil {
			t.Fatalf("could not make request to test server, %v", err)
		}
		defer response.Body.Close()

		got, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("could not read response body, %v", err)
		}

		want := "Not Found\n"

		assertResponseBody(t, string(got), want)
		assertResponseCode(t, response.StatusCode, http.StatusNotFound)

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

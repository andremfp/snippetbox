package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andremfp/snippetbox/internal/database"
)

type StubSnippetStore struct {
	Snippets []database.Snippet
}

func (s *StubSnippetStore) Insert(title, content string, expires int) (int, error) {
	testSnippet := database.Snippet{
		ID:      1,
		Title:   "test title",
		Content: "test content",
		Created: time.Date(2024, time.March, 21, 16, 17, 51, 0, time.UTC),
		Expires: time.Date(2024, time.March, 21, 17, 17, 51, 0, time.UTC),
	}

	s.Snippets = append(s.Snippets, testSnippet)

	return testSnippet.ID, nil
}

func (s *StubSnippetStore) Get(id int) (*database.Snippet, error) {
	if id != 1 {
		return nil, database.ErrNoRecord
	}
	return &s.Snippets[0], nil
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

		wantSnippet := "&{ID:1 Title:test title Content:test content Created:2024-03-21 16:17:51 +0000 UTC Expires:2024-03-21 17:17:51 +0000 UTC}"

		id := 1
		// Create an entry
		createResponse, err := testClient.Post(fmt.Sprintf("%s/snippet/create", testServer.URL), "", nil)
		if err != nil {
			t.Fatalf("could not make request to test server, %v", err)
		}
		defer createResponse.Body.Close()

		// Get the snippet created previously
		getResponse, err := testClient.Get(fmt.Sprintf("%s/snippet/view?id=%d", testServer.URL, id))
		if err != nil {
			t.Fatalf("could not make request to test server, %v", err)
		}
		defer getResponse.Body.Close()

		gotSnippet, err := io.ReadAll(getResponse.Body)
		if err != nil {
			t.Fatalf("could not read response body, %v", err)
		}

		if string(gotSnippet) != wantSnippet {
			t.Errorf("got snippet %s, want %s", gotSnippet, wantSnippet)
		}
		assertResponseCode(t, getResponse.StatusCode, http.StatusOK)

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

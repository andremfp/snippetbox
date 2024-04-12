package server_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/andremfp/snippetbox/internal/database"
	"github.com/andremfp/snippetbox/internal/server"
	"github.com/andremfp/snippetbox/internal/templates"
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
	if id > len(s.Snippets) {
		return nil, database.ErrNoRecord
	}

	return &s.Snippets[0], nil
}

func (s *StubSnippetStore) Latest() ([]*database.Snippet, error) {

	return nil, nil
}

var testApp = &server.Application{
	InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
}

func TestServer(t *testing.T) {

	testApp.SnippetStore = &StubSnippetStore{}
	testServer := httptest.NewServer(testApp.NewServeMux())
	testClient := testServer.Client()
	testClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		t.Errorf("failed to create template cache: %v", err)
	}

	testApp.TemplateCache = templateCache

	defer testServer.Close()

	t.Run("root path returns 200", func(t *testing.T) {

		response, err := testClient.Get(fmt.Sprintf("%s/", testServer.URL))
		if err != nil {
			t.Fatalf("could not make request to test server, %v", err)
		}

		assertResponseCode(t, response.StatusCode, http.StatusOK)

	})

	t.Run("display existing snippet returns 200", func(t *testing.T) {

		// Create a snippet
		_, err := testClient.Post(fmt.Sprintf("%s/snippet/create", testServer.URL), "", nil)
		if err != nil {
			t.Fatalf("could not make create request to test server, %v", err)
		}

		id := 1
		// Get the snippet created previously
		getResponse, err := testClient.Get(fmt.Sprintf("%s/snippet/view?id=%d", testServer.URL, id))
		if err != nil {
			t.Fatalf("could not make get request to test server, %v", err)
		}

		assertResponseCode(t, getResponse.StatusCode, http.StatusOK)

	})

	t.Run("snippet not found", func(t *testing.T) {

		// Get a snippet that does not exist
		id := 2
		response, err := testClient.Get(fmt.Sprintf("%s/snippet/view?id=%d", testServer.URL, id))
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

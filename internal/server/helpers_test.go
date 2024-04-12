package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andremfp/snippetbox/internal/database"
	"github.com/andremfp/snippetbox/internal/templates"
	approvals "github.com/approvals/go-approval-tests"
)

var testSnippets = []*database.Snippet{
	{
		ID:      1,
		Title:   "title1",
		Content: "content1",
		Created: time.Date(2024, time.March, 21, 16, 17, 51, 0, time.UTC),
		Expires: time.Date(2024, time.March, 21, 17, 17, 51, 0, time.UTC),
	},
	{
		ID:      2,
		Title:   "title2",
		Content: "content2",
		Created: time.Date(2024, time.March, 21, 16, 17, 51, 0, time.UTC),
		Expires: time.Date(2024, time.March, 22, 17, 17, 51, 0, time.UTC),
	},
}

func TestRender(t *testing.T) {

	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		t.Errorf("failed to create template cache: %v", err)
	}

	testApp.TemplateCache = templateCache

	tests := []struct {
		name           string
		templateName   string
		data           *templates.TemplateData
		expectedStatus int
	}{
		{
			name:         "home page is rendered successfully and valid",
			templateName: "home.html",
			data:         &templates.TemplateData{CurrentYear: time.Now().Year(), Snippets: testSnippets},
		},
		{
			name:         "view page is rendered successfully and valid",
			templateName: "view.html",
			data:         &templates.TemplateData{CurrentYear: time.Now().Year(), Snippet: testSnippets[0]},
		},
	}

	for _, tt := range tests {

		t.Run("home page is rendered successfully and valid", func(t *testing.T) {
			w := httptest.NewRecorder()

			testApp.Render(w, http.StatusOK, tt.templateName, tt.data)

			approvals.VerifyString(t, w.Body.String())

		})

	}
}

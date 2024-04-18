package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/andremfp/snippetbox/internal/database"
	"github.com/andremfp/snippetbox/internal/server"
	"github.com/andremfp/snippetbox/internal/templates"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/go-playground/form/v4"
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

func TestDecodePostForm(t *testing.T) {
	type testDestinationForm struct {
		Key1 string `form:"key1"`
		Key2 string `form:"key2"`
	}

	var dst testDestinationForm

	tests := []struct {
		name      string
		formData  string
		expectErr bool
		dst       *testDestinationForm
	}{
		{
			name:      "Success case",
			formData:  "key1=value1&key2=value2",
			dst:       &dst,
			expectErr: false,
		},
		{
			name:      "Parsing form failure",
			formData:  "",
			dst:       nil,
			expectErr: true,
		},

		{
			name:      "Decoding form data failure",
			formData:  "key1=value1&key2=value2",
			dst:       nil,
			expectErr: true,
		},
	}

	app := &server.Application{
		FormDecoder: form.NewDecoder(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.formData == "" {
				req, err := http.NewRequest("POST", "/example", nil)
				if err != nil {
					t.Fatal(err)
				}

				err = app.DecodePostForm(req, tt.dst)

				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.formData))

				defer func() {
					if r := recover(); r != nil {
						t.Logf("Panic occurred as expected: %v", r)
					} else {
						if tt.dst == nil {
							t.Errorf("Expected panic, but got none")
						}
					}
				}()

				err := app.DecodePostForm(req, tt.dst)

				if err != nil {
					t.Errorf("Expected no error, but got %v", err)
				}
			}

		})
	}
}

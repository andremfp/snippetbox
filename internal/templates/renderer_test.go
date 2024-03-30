package templates_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/andremfp/snippetbox/internal/database"
	"github.com/andremfp/snippetbox/internal/templates"
	approvals "github.com/approvals/go-approval-tests"
)

func TestRenderTemplate(t *testing.T) {
	t.Run("home page is rendered successfully and valid", func(t *testing.T) {
		buf := bytes.Buffer{}

		htmlFiles := []string{
			"ui/html/base.html",
			"ui/html/partials/nav.html",
			"ui/html/pages/home.html",
		}

		if err := templates.RenderTemplate(&buf, htmlFiles, nil); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())

	})

	t.Run("view page is rendered successfully and valid", func(t *testing.T) {
		buf := bytes.Buffer{}

		testSnippet := &database.Snippet{
			ID:      1,
			Title:   "title",
			Content: "content",
			Created: time.Date(2024, time.March, 21, 16, 17, 51, 0, time.UTC),
			Expires: time.Date(2024, time.March, 21, 17, 17, 51, 0, time.UTC),
		}

		htmlFiles := []string{
			"ui/html/base.html",
			"ui/html/partials/nav.html",
			"ui/html/pages/view.html",
		}

		if err := templates.RenderTemplate(&buf, htmlFiles, testSnippet); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())

	})
}

package html_test

import (
	"bytes"
	"testing"

	"github.com/andremfp/snippetbox/internal/html"
	approvals "github.com/approvals/go-approval-tests"
)

func TestRenderTemplate(t *testing.T) {
	t.Run("home page is rendered successfully and valid", func(t *testing.T) {
		buf := bytes.Buffer{}

		htmlFiles := []string{
			"../../ui/html/base.html",
			"../../ui/html/partials/nav.html",
			"../../ui/html/pages/home.html",
		}

		if err := html.RenderTemplate(&buf, htmlFiles); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())

	})
}

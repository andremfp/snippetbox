package templates_test

import (
	"testing"

	"github.com/andremfp/snippetbox/internal/templates"
)

func TestNewTemplateCache(t *testing.T) {

	numPages := 3

	want := []string{
		"home.html",
		"view.html",
		"create.html",
	}

	cache, err := templates.NewTemplateCache()
	if err != nil {
		t.Errorf("failed to create template cache: %v", err)
	}

	if len(cache) != numPages {
		t.Errorf("want template cache with %d pages, got %d", numPages, len(cache))
	}

	for _, page := range want {
		if _, ok := cache[page]; !ok {
			t.Errorf("want template %s in template cache", page)
		}
	}
}

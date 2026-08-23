package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func get(t *testing.T, path string) (int, string) {
	t.Helper()
	req := httptest.NewRequest("GET", path, nil)
	rec := httptest.NewRecorder()
	handler().ServeHTTP(rec, req)
	resp := rec.Result()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rec.Code, string(body)
}

func TestServesIndex(t *testing.T) {
	code, body := get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, "Situation") || !strings.Contains(body, "Behavior") || !strings.Contains(body, "Impact") {
		t.Error("index.html is missing the three SBI beats")
	}
	if !strings.Contains(body, "kindel.com") {
		t.Error("index.html is missing the required visible link to kindel.com")
	}
}

func TestServesAssets(t *testing.T) {
	for _, path := range []string{"/style.css", "/app.js"} {
		if code, _ := get(t, path); code != 200 {
			t.Errorf("GET %s = %d, want 200", path, code)
		}
	}
}

func TestNoAbsoluteAssetPaths(t *testing.T) {
	// The app must work when hosted under a prefix like /apps/sbi/, so
	// index.html may not reference assets from the site root.
	_, body := get(t, "/")
	for _, bad := range []string{`href="/style`, `src="/app`, `href="/icon`} {
		if strings.Contains(body, bad) {
			t.Errorf("index.html contains root-absolute asset path %q", bad)
		}
	}
}

func TestIndexNotMarketedAsTool(t *testing.T) {
	_, body := get(t, "/")
	lower := strings.ToLower(body)
	if strings.Contains(lower, "tool") && !strings.Contains(lower, "tooltip") {
		t.Error("index.html markets SBI as a 'tool' instead of an 'app'")
	}
}

func TestReadmeNotMarketedAsTool(t *testing.T) {
	data, err := os.ReadFile("README.md")
	if err != nil {
		t.Skip("README.md not found")
	}
	lower := strings.ToLower(string(data))
	// Exclude legitimate uses: tooltip, tooling, and /tools/ path references.
	if strings.Contains(lower, "tool") &&
		!strings.Contains(lower, "tooltip") &&
		!strings.Contains(lower, "tooling") &&
		!strings.Contains(lower, "/tools/") {
		t.Error("README.md markets SBI as a 'tool' instead of an 'app'")
	}
}

func TestIndexHasKindelLink(t *testing.T) {
	_, body := get(t, "/")
	if !strings.Contains(body, "kindel.com") {
		t.Error("index.html missing required visible link to kindel.com")
	}
}

func TestIndexHasSBIWords(t *testing.T) {
	_, body := get(t, "/")
	for _, word := range []string{"Situation", "Behavior", "Impact"} {
		if !strings.Contains(body, word) {
			t.Errorf("index.html missing SBI word: %s", word)
		}
	}
}

func TestCardJSONHref(t *testing.T) {
	data, err := os.ReadFile("card.json")
	if err != nil {
		t.Skip("card.json not found")
	}
	var card struct {
		Href string `json:"href"`
	}
	if err := json.Unmarshal(data, &card); err != nil {
		t.Fatalf("card.json parse error: %v", err)
	}
	if card.Href != "/kld/apps/sbi/" {
		t.Errorf("card.json href = %q, want /kld/apps/sbi/", card.Href)
	}
}

func TestTextareasHaveAccessibleNames(t *testing.T) {
	_, body := get(t, "/")
	textareaIDs := []string{"draft", "situation", "behavior", "impact"}
	for _, id := range textareaIDs {
		hasAriaLabel := strings.Contains(body, `id="`+id+`"`) &&
			(strings.Contains(body, `<textarea id="`+id+`"`) || strings.Contains(body, `<textarea`))
		if !hasAriaLabel {
			t.Errorf("textarea %q not found in index.html", id)
			continue
		}
		// Check for aria-label, aria-labelledby, or <label for="id">
		hasAccessibleName := false
		// Simple check: find the textarea tag and verify it has aria-label
		if strings.Contains(body, `id="`+id+`"`) {
			// Look for aria-label on the same element
			if strings.Contains(body, `<textarea id="`+id+`"`) {
				// Check if the line containing this textarea has aria-label
				// Using a simple pattern match
				pattern := `id="` + id + `"`
				idx := strings.Index(body, pattern)
				if idx != -1 {
					// Find the enclosing tag
					start := strings.LastIndex(body[:idx], "<textarea")
					if start != -1 {
						end := strings.Index(body[start:], ">")
						if end != -1 {
							tag := body[start : start+end+1]
							if strings.Contains(tag, "aria-label=") {
								hasAccessibleName = true
							}
						}
					}
				}
			}
			// Check for <label for="id">
			if strings.Contains(body, `<label for="`+id+`"`) ||
				strings.Contains(body, `<label for='`+id+`'`) {
				hasAccessibleName = true
			}
			// Check for aria-labelledby on the textarea
			if strings.Contains(body, `id="`+id+`"`) && strings.Contains(body, `aria-labelledby=`) {
				pattern := `id="` + id + `"`
				idx := strings.Index(body, pattern)
				if idx != -1 {
					start := strings.LastIndex(body[:idx], "<textarea")
					if start != -1 {
						end := strings.Index(body[start:], ">")
						if end != -1 {
							tag := body[start : start+end+1]
							if strings.Contains(tag, "aria-labelledby=") {
								hasAccessibleName = true
							}
						}
					}
				}
			}
		}
		if !hasAccessibleName {
			t.Errorf("textarea %q missing accessible name (aria-label, aria-labelledby, or <label for>)", id)
		}
	}
}

func TestNoEmDashesInUserFacingEnglish(t *testing.T) {
	files := map[string]func() ([]byte, error){
		"index.html": func() ([]byte, error) {
			_, body := get(t, "/")
			return []byte(body), nil
		},
		"app.js": func() ([]byte, error) {
			_, body := get(t, "/app.js")
			return []byte(body), nil
		},
		"README.md": func() ([]byte, error) {
			return os.ReadFile("README.md")
		},
	}
	for name, loader := range files {
		data, err := loader()
		if err != nil {
			t.Errorf("cannot read %s: %v", name, err)
			continue
		}
		content := string(data)
		// Check for U+2014 (em dash) character
		if strings.Contains(content, "\u2014") {
			t.Errorf("%s contains em dash (U+2014)", name)
		}
		// Check for &mdash; HTML entity
		if strings.Contains(content, "&mdash;") {
			t.Errorf("%s contains &mdash; HTML entity", name)
		}
	}
}

func TestLintboxScrollbarGeometry(t *testing.T) {
	data, err := os.ReadFile("static/style.css")
	if err != nil {
		t.Fatal("cannot read static/style.css:", err)
	}
	css := string(data)

	// We need to verify that .lintbox .backdrop and .lintbox textarea have
	// matching overflow-y and scrollbar-gutter settings.
	//
	// The fix should have a shared rule for both (or matching individual rules).
	// We check that:
	// 1. Both have overflow-y: auto (not hidden-only on backdrop)
	// 2. scrollbar-gutter: stable is present

	// Look for the shared rule ".lintbox .backdrop,\n.lintbox textarea"
	hasSharedRule := strings.Contains(css, ".lintbox .backdrop,") &&
		strings.Contains(css, ".lintbox textarea")

	// Check for overflow-y: auto in the shared context
	hasOverflowYAuto := strings.Contains(css, "overflow-y: auto")

	// Check for scrollbar-gutter: stable
	hasScrollbarGutter := strings.Contains(css, "scrollbar-gutter: stable")

	// Verify backdrop does not have overflow: hidden without matching gutter
	// Find backdrop-specific rules
	backdropIdx := strings.Index(css, ".lintbox .backdrop {")
	if backdropIdx != -1 {
		// Find the end of this rule block
		end := strings.Index(css[backdropIdx:], "}")
		if end != -1 {
			backdropRule := css[backdropIdx : backdropIdx+end]
			// If backdrop has overflow: hidden, it must be in a separate rule
			// and the shared rule must handle scrolling
			if strings.Contains(backdropRule, "overflow: hidden") &&
				!strings.Contains(backdropRule, "scrollbar-gutter") {
				// This is only a problem if there's no shared rule with proper settings
				if !hasSharedRule || !hasScrollbarGutter {
					t.Error("backdrop uses overflow: hidden without shared scrollbar-gutter rule")
				}
			}
		}
	}

	if !hasOverflowYAuto {
		t.Error("style.css missing overflow-y: auto for lintbox textarea/backdrop")
	}

	if !hasScrollbarGutter {
		t.Error("style.css missing scrollbar-gutter: stable for lintbox textarea/backdrop")
	}

	if !hasSharedRule {
		t.Error("style.css should have shared rule for .lintbox .backdrop and .lintbox textarea")
	}
}

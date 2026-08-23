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
	if card.Href != "/apps/sbi/" {
		t.Errorf("card.json href = %q, want /apps/sbi/", card.Href)
	}
}

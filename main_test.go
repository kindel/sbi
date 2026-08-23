package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func get(t *testing.T, path string) (int, string) {
	t.Helper()
	req := httptest.NewRequest("GET", path, nil)
	rec := httptest.NewRecorder()
	handler().ServeHTTP(rec, req)
	body, err := io.ReadAll(rec.Result().Body)
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
	// The app must work when hosted under /sbi/, so index.html may not
	// reference assets from the site root.
	_, body := get(t, "/")
	for _, bad := range []string{`href="/style`, `src="/app`, `href="/icon`} {
		if strings.Contains(body, bad) {
			t.Errorf("index.html contains root-absolute asset path %q", bad)
		}
	}
}

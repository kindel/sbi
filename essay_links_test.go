package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestEssayLinkStaysOnKindel(t *testing.T) {
	const apex = "https://kindel.com/essays/the-secret-to-giving-feedback-without-triggering-defensiveness/"
	body, err := os.ReadFile("static/index.html")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), apex) {
		t.Fatal("index.html should link the essay on kindel.com")
	}
	card, err := os.ReadFile("card.json")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(card), apex) {
		t.Fatal("card.json should link the essay on kindel.com")
	}
	assertNoEssayPermalinks(t)
}

func assertNoEssayPermalinks(t *testing.T) {
	t.Helper()
	raw, err := os.ReadFile("data/essay_slugs.json")
	if err != nil {
		t.Fatal(err)
	}
	var snap struct {
		BySlug map[string]string `json:"by_slug"`
	}
	if err := json.Unmarshal(raw, &snap); err != nil {
		t.Fatal(err)
	}
	dated := regexp.MustCompile(`https?://(?:www\.)?blog\.kindel\.com/\d{4}/\d{2}/\d{2}/([a-z0-9]+(?:-[a-z0-9]+)*)/?`)
	skip := map[string]bool{
		".git": true, "essay_slugs.json": true, "check-essay-links.js": true,
	}
	var hits []string
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if skip[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if skip[info.Name()] {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".html", ".md", ".js", ".json", ".go", ".css", ".txt":
		default:
			return nil
		}
		text, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, match := range dated.FindAllStringSubmatch(string(text), -1) {
			if _, ok := snap.BySlug[strings.ToLower(match[1])]; ok {
				hits = append(hits, path+": "+match[0])
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) > 0 {
		t.Fatalf("essay permalinks still on the blog:\n%s", strings.Join(hits, "\n"))
	}
}

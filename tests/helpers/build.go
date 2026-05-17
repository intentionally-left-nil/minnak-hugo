// Package helpers provides utilities for Hugo theme tests.
// The central function, Build, runs "hugo" against the exampleSite, and
// optionally runs "npx pagefind" afterwards.  Subsequent calls reuse the
// already-built public directory unless force=true.
package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// repoRoot returns the absolute path of the repository root (one level above
// this file's package directory, i.e. minnak-hugo/).
func RepoRoot() string {
	_, file, _, _ := runtime.Caller(0)
	// file = .../tests/helpers/build.go → go up two dirs
	return filepath.Join(filepath.Dir(file), "..", "..")
}

// ExampleSiteDir returns the absolute path of exampleSite/.
func ExampleSiteDir() string {
	return filepath.Join(RepoRoot(), "exampleSite")
}

// PublicDir returns the absolute path of exampleSite/public/.
func PublicDir() string {
	return filepath.Join(ExampleSiteDir(), "public")
}

// Build runs "hugo" against exampleSite and optionally "npx pagefind"
// afterwards.  It writes build output to t.Log.
// If the public/ directory already exists and force is false the build is
// skipped (fast path for tests running in the same process).
func Build(t *testing.T, force bool, runPagefind bool) {
	t.Helper()

	publicDir := PublicDir()

	if !force {
		if _, err := os.Stat(publicDir); err == nil {
			t.Log("helpers.Build: reusing existing public/ directory")
			return
		}
	}

	exampleSite := ExampleSiteDir()

	t.Log("helpers.Build: running hugo...")
	hugoCmd := exec.Command("hugo", "--source", exampleSite, "--logLevel", "warn")
	hugoCmd.Dir = RepoRoot()
	out, err := hugoCmd.CombinedOutput()
	if err != nil {
		t.Logf("hugo output:\n%s", string(out))
		t.Fatalf("hugo build failed: %v", err)
	}
	t.Logf("hugo build OK (%d bytes output)", len(out))

	if runPagefind {
		t.Log("helpers.Build: running npx pagefind...")
		pfCmd := exec.Command("npx", "--yes", "pagefind@latest", "--site", publicDir)
		pfCmd.Dir = RepoRoot()
		pfOut, pfErr := pfCmd.CombinedOutput()
		if pfErr != nil {
			// Pagefind failure is non-fatal for most markup tests; log a warning.
			t.Logf("WARNING: pagefind failed (search tests will fail): %v\n%s", pfErr, string(pfOut))
		} else {
			t.Logf("pagefind OK (%d bytes output)", len(pfOut))
		}
	}
}

// ParseFile reads a built HTML file from public/ and returns a goquery document.
// path is relative to public/, e.g. "index.html" or "posts/rust-ownership-model/index.html".
func ParseFile(t *testing.T, path string) *goquery.Document {
	t.Helper()
	full := filepath.Join(PublicDir(), filepath.FromSlash(path))
	f, err := os.Open(full)
	if err != nil {
		t.Fatalf("ParseFile: cannot open %s: %v", full, err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Fatalf("ParseFile: cannot parse %s: %v", full, err)
	}
	return doc
}

// AssertSelector asserts that the CSS selector matches exactly `want` elements.
func AssertSelector(t *testing.T, doc *goquery.Document, selector string, want int) {
	t.Helper()
	got := doc.Find(selector).Length()
	if got != want {
		t.Errorf("selector %q: got %d, want %d", selector, got, want)
	}
}

// AssertSelectorAtLeast asserts that the CSS selector matches at least `min` elements.
func AssertSelectorAtLeast(t *testing.T, doc *goquery.Document, selector string, min int) {
	t.Helper()
	got := doc.Find(selector).Length()
	if got < min {
		t.Errorf("selector %q: got %d, want at least %d", selector, got, min)
	}
}

// AssertText asserts that at least one element matching the CSS selector
// contains the given text (substring, case-insensitive).
func AssertText(t *testing.T, doc *goquery.Document, selector, text string) {
	t.Helper()
	found := false
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
		if strings.Contains(strings.ToLower(s.Text()), strings.ToLower(text)) {
			found = true
		}
	})
	if !found {
		t.Errorf("no element matching %q contains text %q", selector, text)
	}
}

// AssertAttr asserts that the first element matching selector has attribute
// attr with value containing substr.
func AssertAttr(t *testing.T, doc *goquery.Document, selector, attr, substr string) {
	t.Helper()
	val, exists := doc.Find(selector).First().Attr(attr)
	if !exists {
		t.Errorf("selector %q: attribute %q not found", selector, attr)
		return
	}
	if !strings.Contains(val, substr) {
		t.Errorf("selector %q attr %q = %q, want to contain %q", selector, attr, val, substr)
	}
}

// FileExists asserts that a file exists at the given public/-relative path.
func FileExists(t *testing.T, path string) {
	t.Helper()
	full := filepath.Join(PublicDir(), filepath.FromSlash(path))
	if _, err := os.Stat(full); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", full)
	}
}

// PostTitles returns the titles of all posts in the exampleSite that should
// appear in the output.  Update if content/ changes.
func PostTitles() []string {
	return []string{
		"Rust's Ownership Model Is Actually About Time",
		"Setting Up a Home Lab on a Decommissioned ThinkPad",
		"What GPT-4 Gets Wrong About Cantonese",
		"Learning Jyutping as an Adult Heritage Speaker",
		"Golden Hour on the Carbon River Road",
		"Winter on the Wonderland Trail: What No One Tells You",
		"Diffusion Models From First Principles",
		"Mount Rainier in Four Frames",
		"Two Years With a Fairphone",
	}
}

// Href extracts the href attribute from the first match of selector.
func Href(doc *goquery.Document, selector string) string {
	val, _ := doc.Find(selector).First().Attr("href")
	return val
}

// TextOf returns the trimmed text of the first match of selector.
func TextOf(doc *goquery.Document, selector string) string {
	return strings.TrimSpace(doc.Find(selector).First().Text())
}

// AllTexts returns the trimmed text of every match of selector.
func AllTexts(doc *goquery.Document, selector string) []string {
	var out []string
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
		out = append(out, strings.TrimSpace(s.Text()))
	})
	return out
}

// AllHrefs returns the href attribute of every match of selector.
func AllHrefs(doc *goquery.Document, selector string) []string {
	var out []string
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
		if h, ok := s.Attr("href"); ok {
			out = append(out, h)
		}
	})
	return out
}

// Fatalf wraps fmt.Sprintf for test messages.
func Fatalf(t *testing.T, format string, args ...any) {
	t.Helper()
	t.Fatalf(format, args...)
}

// Sprintf is a convenience alias.
var Sprintf = fmt.Sprintf

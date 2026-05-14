package tests

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// M3 — baseof.html, dark theme, fonts, fingerprinted CSS
// ----------------------------------------------------------------------------

// TestDarkThemeBodyClass verifies <body class="dark-theme">.
func TestDarkThemeBodyClass(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	bodyClass, _ := doc.Find("body").Attr("class")
	if bodyClass != "dark-theme" {
		t.Errorf("body class: got %q, want %q", bodyClass, "dark-theme")
	}
}

// TestSiteTitle verifies the site title link is present in the header.
func TestSiteTitle(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, "#masthead .site-title a", 1)
	helpers.AssertText(t, doc, "#masthead .site-title a", "terminal.space")
}

// TestCharsetMeta verifies <meta charset="utf-8">.
func TestCharsetMeta(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, `meta[charset="utf-8"]`, 1)
}

// TestViewportMeta verifies the viewport meta tag is present.
func TestViewportMeta(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, `meta[name="viewport"]`, 1)
}

// TestFingerprintedCSS verifies a <link> to a fingerprinted CSS file is present.
func TestFingerprintedCSS(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	found := false
	doc.Find(`link[rel="stylesheet"]`).Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		// Fingerprinted files have a hash in the filename, e.g. index.bundle.abc12345.css
		if len(href) > 20 {
			found = true
		}
	})
	if !found {
		t.Error("expected a fingerprinted stylesheet link")
	}
}

// TestRSSLink verifies an RSS autodiscovery link is present.
func TestRSSLink(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, `link[type="application/rss+xml"]`, 1)
}

// TestFooterCopyright verifies the footer contains the site title.
func TestFooterCopyright(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertText(t, doc, ".site-footer", "terminal.space")
}

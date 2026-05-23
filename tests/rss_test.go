package tests

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// RSS feed
//
// The site is configured to emit a single site-wide feed at /rss.xml (home
// page only).  Tag and category term pages do not get feeds.  Each item in
// the feed contains the full post HTML rather than the truncated summary.
// ----------------------------------------------------------------------------

// rssItem is used to unmarshal <item> elements from rss.xml.
type rssItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
}

// rssFeed is the top-level structure expected in rss.xml.
type rssFeed struct {
	Channel struct {
		Items []rssItem `xml:"item"`
	} `xml:"channel"`
}

// readRSSFeed parses public/rss.xml and returns the structured feed.
// It fails the test immediately if the file is missing or malformed.
func readRSSFeed(t *testing.T) rssFeed {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(helpers.PublicDir(), "rss.xml"))
	if err != nil {
		t.Fatalf("readRSSFeed: %v", err)
	}
	var feed rssFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("readRSSFeed: parse error: %v", err)
	}
	return feed
}

// TestRSSFeedExistsAtRSSXML verifies the site-wide feed is published to
// /rss.xml (baseName = "rss" in [outputFormats]).
func TestRSSFeedExistsAtRSSXML(t *testing.T) {
	buildOnce(t)
	helpers.FileExists(t, "rss.xml")
}

// TestRSSIndexXMLAbsent verifies that /index.xml — Hugo's default RSS path —
// is not emitted now that baseName = "rss".
func TestRSSIndexXMLAbsent(t *testing.T) {
	buildOnce(t)
	full := filepath.Join(helpers.PublicDir(), "index.xml")
	if _, err := os.Stat(full); err == nil {
		t.Error("unexpected file index.xml: feed should be at rss.xml, not index.xml")
	}
}

// TestTagPagesHaveNoRSSFeed verifies individual tag term pages do not emit
// their own feeds (outputs.term = ['html'] only).
func TestTagPagesHaveNoRSSFeed(t *testing.T) {
	buildOnce(t)
	for _, slug := range []string{"rust", "ml", "hiking"} {
		path := filepath.Join(helpers.PublicDir(), "tag", slug, "index.xml")
		if _, err := os.Stat(path); err == nil {
			t.Errorf("unexpected RSS feed at tag/%s/index.xml", slug)
		}
	}
}

// TestCategoryPagesHaveNoRSSFeed verifies category term pages do not emit
// their own feeds (outputs.term = ['html'] only).
func TestCategoryPagesHaveNoRSSFeed(t *testing.T) {
	buildOnce(t)
	for _, slug := range []string{"technology", "ai", "photography"} {
		path := filepath.Join(helpers.PublicDir(), "category", slug, "index.xml")
		if _, err := os.Stat(path); err == nil {
			t.Errorf("unexpected RSS feed at category/%s/index.xml", slug)
		}
	}
}

// TestRSSAutodiscoveryLinkPointsToRSSXML verifies the <link> autodiscovery
// tag in <head> references /rss.xml, not /index.xml.
func TestRSSAutodiscoveryLinkPointsToRSSXML(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertAttr(t, doc, `link[type="application/rss+xml"]`, "href", "rss.xml")
}

// TestRSSAutodiscoveryAbsentOnTaxonomyPages verifies that tag and category
// pages carry no RSS autodiscovery link, since RSS output is restricted to
// the home page.
func TestRSSAutodiscoveryAbsentOnTaxonomyPages(t *testing.T) {
	buildOnce(t)
	for _, path := range []string{
		"tag/rust/index.html",
		"category/technology/index.html",
	} {
		doc := helpers.ParseFile(t, path)
		got := doc.Find(`link[type="application/rss+xml"]`).Length()
		if got != 0 {
			t.Errorf("%s: expected no RSS autodiscovery link, got %d", path, got)
		}
	}
}

// TestRSSFeedContainsAllPosts verifies every post title from the exampleSite
// appears as an <item> in the feed.
func TestRSSFeedContainsAllPosts(t *testing.T) {
	buildOnce(t)
	feed := readRSSFeed(t)

	inFeed := make(map[string]bool, len(feed.Channel.Items))
	for _, item := range feed.Channel.Items {
		inFeed[item.Title] = true
	}

	for _, title := range helpers.PostTitles() {
		if !inFeed[title] {
			t.Errorf("RSS feed missing post %q", title)
		}
	}
}

// TestRSSFeedDescriptionIsFullContent verifies every item's <description>
// contains full post HTML rather than the plain-text truncated summary.
// The description must include at least one HTML tag (e.g. <p>).
func TestRSSFeedDescriptionIsFullContent(t *testing.T) {
	buildOnce(t)
	feed := readRSSFeed(t)

	if len(feed.Channel.Items) == 0 {
		t.Fatal("RSS feed has no items")
	}

	for _, item := range feed.Channel.Items {
		// encoding/xml unescapes &lt; → <, so a description with HTML will
		// contain "<" after parsing.
		if !strings.Contains(item.Description, "<") {
			t.Errorf("item %q: description appears to be plain-text summary, expected full HTML content", item.Title)
		}
	}
}

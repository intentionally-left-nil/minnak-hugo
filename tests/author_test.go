package tests

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// Guest author
//
// Posts may declare an optional `author` in front matter, as either a plain
// string (`author: "Jane Doe"`) or a map (`{name, url, email}`). When set, a
// "Guest author: NAME" byline appears under the post title and SEO meta
// (meta[name=author], og article:author, JSON-LD Article) is emitted in
// <head>. Posts without `author` render unchanged.
//
// Fixtures (exampleSite/content/posts/):
//   jyutping-heritage-speaker.md  - string form: "Maya Chen"
//   gpt4-cantonese/               - map form:    "Dr. Wing Lam" + url
//   rust-ownership-model/         - no author (baseline)
// ----------------------------------------------------------------------------

const (
	authorStringPostPath = "posts/jyutping-heritage-speaker/index.html"
	authorMapPostPath    = "posts/gpt4-cantonese/index.html"
	authorAbsentPostPath = "posts/rust-ownership-model/index.html"
)

// TestBylineStringAuthor verifies a string-form author renders as plain text.
func TestBylineStringAuthor(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, authorStringPostPath)

	helpers.AssertSelector(t, doc, ".entry-header .byline", 1)
	helpers.AssertText(t, doc, ".byline", "Guest author")
	helpers.AssertText(t, doc, ".byline", "Maya Chen")

	// String form has no url, so byline text contains the name without an <a>.
	if doc.Find(".byline a").Length() != 0 {
		t.Errorf("string-form author byline should not contain an <a>")
	}
}

// TestBylineMapAuthorLinked verifies a map-form author with `url` renders as
// a link wrapping the name.
func TestBylineMapAuthorLinked(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, authorMapPostPath)

	helpers.AssertSelector(t, doc, ".entry-header .byline", 1)
	helpers.AssertText(t, doc, ".byline", "Guest author")
	helpers.AssertText(t, doc, ".byline a", "Dr. Wing Lam")
	helpers.AssertAttr(t, doc, ".byline a", "href", "https://example.com/wing-lam")
	helpers.AssertAttr(t, doc, ".byline a", "rel", "author")
}

// TestBylineAbsentWithoutAuthor verifies posts without front-matter author
// render no byline, no meta[name=author], and no JSON-LD Article block.
func TestBylineAbsentWithoutAuthor(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, authorAbsentPostPath)

	helpers.AssertSelector(t, doc, ".byline", 0)
	helpers.AssertSelector(t, doc, `meta[name="author"]`, 0)
	helpers.AssertSelector(t, doc, `meta[property="article:author"]`, 0)
	helpers.AssertSelector(t, doc, `script[type="application/ld+json"]`, 0)
}

// TestArticleJSONLDOnlyOnPosts verifies the Article JSON-LD is not emitted
// on standalone pages (type: page) even when search engines might otherwise
// pick it up. The home page is also covered as a kind != "page" case.
func TestArticleJSONLDOnlyOnPosts(t *testing.T) {
	buildOnce(t)
	for _, path := range []string{
		"about/index.html",
		"index.html",
		"category/technology/index.html",
	} {
		doc := helpers.ParseFile(t, path)
		helpers.AssertSelector(t, doc, `script[type="application/ld+json"]`, 0)
		helpers.AssertSelector(t, doc, `meta[name="author"]`, 0)
	}
}

// TestRSSGuestAuthorOmitsAuthorTag verifies guest-author posts whose front
// matter has no email do not get an RSS <author> element. RSS spec requires
// <author> to be an email; misattributing the post to the site owner's
// email would be incorrect, so the element is omitted entirely.
//
// Conversely, when a guest author supplies an email, the <author> element
// is emitted with that email and the author's name in parens, matching the
// RSS 2.0 format `email (Name)`.
func TestRSSGuestAuthorOmitsAuthorTag(t *testing.T) {
	buildOnce(t)

	type rssItem struct {
		Title  string `xml:"title"`
		Link   string `xml:"link"`
		Author string `xml:"author"`
	}
	type rssFeed struct {
		Channel struct {
			Items []rssItem `xml:"item"`
		} `xml:"channel"`
	}

	data, err := os.ReadFile(filepath.Join(helpers.PublicDir(), "rss.xml"))
	if err != nil {
		t.Fatalf("read rss.xml: %v", err)
	}
	var feed rssFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("parse rss.xml: %v", err)
	}

	// Expected per-item <author> values keyed by post title.
	// Empty string means "<author> element should be absent".
	want := map[string]string{
		"Learning Jyutping as an Adult Heritage Speaker": "",
		"What GPT-4 Gets Wrong About Cantonese":          "winglam@example.com (Dr. Wing Lam)",
	}
	seen := 0
	for _, item := range feed.Channel.Items {
		expected, tracked := want[item.Title]
		if !tracked {
			continue
		}
		seen++
		if item.Author != expected {
			t.Errorf("item %q: <author>: got %q, want %q", item.Title, item.Author, expected)
		}
	}
	if seen != len(want) {
		t.Errorf("found %d tracked items in feed, want %d", seen, len(want))
	}
}

// TestSEOMetaAuthorPresent verifies meta[name=author] and og article:author
// are emitted on posts with an author. Both string and map author forms are
// covered to confirm both paths through partials/author.html.
func TestSEOMetaAuthorPresent(t *testing.T) {
	buildOnce(t)

	cases := []struct {
		path string
		name string
	}{
		{authorStringPostPath, "Maya Chen"},
		{authorMapPostPath, "Dr. Wing Lam"},
	}
	for _, tc := range cases {
		doc := helpers.ParseFile(t, tc.path)
		helpers.AssertAttr(t, doc, `meta[name="author"]`, "content", tc.name)
		helpers.AssertAttr(t, doc, `meta[property="article:author"]`, "content", tc.name)
	}
}

// TestJSONLDArticlePresent verifies the JSON-LD Article block on a post with
// a string author has the expected schema.org Person author.
func TestJSONLDArticlePresent(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, authorStringPostPath)

	scripts := doc.Find(`script[type="application/ld+json"]`)
	if scripts.Length() != 1 {
		t.Fatalf("expected 1 JSON-LD <script>, got %d", scripts.Length())
	}

	var parsed map[string]any
	raw := strings.TrimSpace(scripts.First().Text())
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		t.Fatalf("JSON-LD failed to parse: %v\nraw: %s", err, raw)
	}

	if got := parsed["@type"]; got != "Article" {
		t.Errorf("@type: got %v, want %q", got, "Article")
	}

	author, ok := parsed["author"].(map[string]any)
	if !ok {
		t.Fatalf("author: expected object, got %T", parsed["author"])
	}
	if got := author["@type"]; got != "Person" {
		t.Errorf("author.@type: got %v, want %q", got, "Person")
	}
	if got := author["name"]; got != "Maya Chen" {
		t.Errorf("author.name: got %v, want %q", got, "Maya Chen")
	}
	// String-form author has no url; absence is the expected behaviour.
	if _, hasURL := author["url"]; hasURL {
		t.Errorf("author.url: should be omitted for string-form author")
	}
}

// TestJSONLDArticleAuthorURL verifies map-form authors include a url in the
// JSON-LD author object.
func TestJSONLDArticleAuthorURL(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, authorMapPostPath)

	scripts := doc.Find(`script[type="application/ld+json"]`)
	if scripts.Length() != 1 {
		t.Fatalf("expected 1 JSON-LD <script>, got %d", scripts.Length())
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(scripts.First().Text())), &parsed); err != nil {
		t.Fatalf("JSON-LD failed to parse: %v", err)
	}

	author, _ := parsed["author"].(map[string]any)
	if got := author["name"]; got != "Dr. Wing Lam" {
		t.Errorf("author.name: got %v, want %q", got, "Dr. Wing Lam")
	}
	if got := author["url"]; got != "https://example.com/wing-lam" {
		t.Errorf("author.url: got %v, want %q", got, "https://example.com/wing-lam")
	}
}

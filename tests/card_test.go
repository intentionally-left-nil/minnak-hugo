package tests

import (
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// M4 — Card partial
// ----------------------------------------------------------------------------

// TestCardStructureExists verifies cards are rendered on the homepage.
func TestCardStructureExists(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelectorAtLeast(t, doc, ".gh-card", 1)
}

// TestCardHasTitle verifies each card has a title link.
func TestCardHasTitle(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelectorAtLeast(t, doc, ".gh-card .gh-card-title a[href]", 1)
}

// TestCardHasDate verifies each card has a <time> element.
func TestCardHasDate(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelectorAtLeast(t, doc, ".gh-card .gh-card-meta time", 1)
}

// TestCardHasReadMore verifies each card has a "Read more" link.
func TestCardHasReadMore(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelectorAtLeast(t, doc, ".gh-card .read-more-link[href]", 1)
}

// TestCardCategoryLink verifies that cards with a category show a tag link.
func TestCardCategoryLink(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	// At least one card should have a .gh-card-tag link
	helpers.AssertSelectorAtLeast(t, doc, ".gh-card .gh-card-tag[href]", 1)
}

// TestCardCategoryLinkPoints verifies a category link points to /category/.
func TestCardCategoryLinkPoints(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	href := helpers.Href(doc, ".gh-card .gh-card-tag")
	if href == "" {
		t.Error("expected a category link on a card")
		return
	}
	const want = "/category/"
	if len(href) < len(want) || href[:len(want)] != want {
		t.Errorf("card category link: got %q, want prefix %q", href, want)
	}
}

// TestAllPostTitlesOnHomepage verifies all 7 fixture post titles appear somewhere
// in the homepage HTML (across the 3-col/2-col/1-col grid copies).
func TestAllPostTitlesOnHomepage(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	titles := helpers.AllTexts(doc, ".gh-card-title a")
	titleSet := make(map[string]bool, len(titles))
	for _, tt := range titles {
		titleSet[tt] = true
	}

	for _, want := range helpers.PostTitles() {
		if !titleSet[want] {
			t.Errorf("post title %q not found on homepage (got: %v)", want, titles)
		}
	}
}

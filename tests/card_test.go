package tests

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
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

// ----------------------------------------------------------------------------
// Card excerpt behaviour (matches WordPress/Ghost 35-word plain-text truncation)
// ----------------------------------------------------------------------------

// excerptCard finds the card for the fixture post "Card Excerpt Test Fixture"
// and returns the text of its excerpt <p>.
func excerptCard(t *testing.T) string {
	t.Helper()
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	var found string
	doc.Find(".gh-card").Each(func(_ int, card *goquery.Selection) {
		if strings.Contains(card.Find(".gh-card-title").Text(), "Card Excerpt Test Fixture") {
			found = card.Find(".gh-card-content p").First().Text()
		}
	})
	if found == "" {
		t.Fatal("card for 'Card Excerpt Test Fixture' not found on homepage")
	}
	return found
}

// TestCardExcerptWordCount verifies the excerpt is at most 35 words.
func TestCardExcerptWordCount(t *testing.T) {
	text := excerptCard(t)
	// Strip trailing " ..." before counting.
	trimmed := strings.TrimSuffix(text, " ...")
	words := strings.Fields(trimmed)
	if len(words) > 35 {
		t.Errorf("excerpt has %d words, want ≤ 35\nexcerpt: %q", len(words), text)
	}
}

// TestCardExcerptTruncated verifies long posts end with " ...".
func TestCardExcerptTruncated(t *testing.T) {
	text := excerptCard(t)
	if !strings.HasSuffix(text, " ...") {
		t.Errorf("excerpt for long post does not end with \" ...\"\nexcerpt: %q", text)
	}
}

// TestCardExcerptNoImages verifies no <img> or <figure> tags appear in any card excerpt.
func TestCardExcerptNoImages(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	doc.Find(".gh-card-content p").Each(func(_ int, p *goquery.Selection) {
		if p.Find("img").Length() > 0 {
			t.Errorf("card excerpt contains an <img> tag: %q", p.Text())
		}
		if p.Find("figure").Length() > 0 {
			t.Errorf("card excerpt contains a <figure> tag: %q", p.Text())
		}
	})
}

// TestCardExcerptNoHeadings verifies no heading tags appear in any card excerpt.
func TestCardExcerptNoHeadings(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	for _, sel := range []string{"h1", "h2", "h3", "h4", "h5", "h6"} {
		doc.Find(".gh-card-content p").Each(func(_ int, p *goquery.Selection) {
			if p.Find(sel).Length() > 0 {
				t.Errorf("card excerpt contains a <%s> tag: %q", sel, p.Text())
			}
		})
	}
}

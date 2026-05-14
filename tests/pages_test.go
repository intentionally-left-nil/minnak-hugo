package tests

import (
	"strings"
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// M7 — Single post page
// ----------------------------------------------------------------------------

const rustPostPath = "posts/rust-ownership-model/index.html"

// TestPostTitle verifies the post title is rendered in the h1.
func TestPostTitle(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, rustPostPath)
	helpers.AssertText(t, doc, ".entry-title", "Rust's Ownership Model Is Actually About Time")
}

// TestPostDateFormat verifies the date is rendered in "Month Day, Year" format.
func TestPostDateFormat(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, rustPostPath)

	timeEl := doc.Find(".posted-on time")
	if timeEl.Length() == 0 {
		t.Fatal("expected .posted-on time element")
	}

	text := strings.TrimSpace(timeEl.Text())
	// Expect "March 15, 2026"
	if !strings.Contains(text, "2026") || !strings.Contains(text, "March") {
		t.Errorf("date text: got %q, want 'March 15, 2026'", text)
	}

	// datetime attr should be machine-readable
	dt, _ := timeEl.Attr("datetime")
	if dt != "2026-03-15" {
		t.Errorf("datetime attr: got %q, want %q", dt, "2026-03-15")
	}
}

// TestPostCategoryLink verifies the primary category link is rendered.
func TestPostCategoryLink(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, rustPostPath)

	catLink := doc.Find(".cat-links a")
	if catLink.Length() == 0 {
		t.Fatal("expected .cat-links a on post page")
	}

	text := strings.TrimSpace(catLink.Text())
	if text != "Technology" {
		t.Errorf("category link text: got %q, want %q", text, "Technology")
	}

	href := helpers.Href(doc, ".cat-links a")
	if !strings.Contains(href, "technology") {
		t.Errorf("category link href: got %q, want to contain 'technology'", href)
	}
}

// TestPostContentRendered verifies the post body markdown is rendered to HTML.
func TestPostContentRendered(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, rustPostPath)
	// The post body contains an h2 "The real insight"
	helpers.AssertText(t, doc, ".entry-content h2", "The real insight")
}

// TestPostTagsInFooter verifies the tag pills are rendered in entry-footer.
func TestPostTagsInFooter(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, rustPostPath)

	// Rust post has tags: rust, programming, memory
	helpers.AssertSelectorAtLeast(t, doc, ".entry-footer .tags-links a", 1)
	helpers.AssertText(t, doc, ".entry-footer .tags-links", "rust")
}

// TestPostPrevNextNav verifies prev/next navigation is rendered on a post that
// has both neighbours.
func TestPostPrevNextNav(t *testing.T) {
	buildOnce(t)
	// The "Diffusion Models" post (2026-03-10) should have both next and prev.
	doc := helpers.ParseFile(t, "posts/diffusion-models-first-principles/index.html")

	// With Hugo's default ordering (descending date), PrevInSection is newer,
	// NextInSection is older.
	helpers.AssertSelector(t, doc, ".gh-readmore", 1)
	helpers.AssertSelector(t, doc, ".gh-readmore-inner", 1)
}

// ----------------------------------------------------------------------------
// M8 — Single page (page/single.html)
// ----------------------------------------------------------------------------

// TestAboutPageTitle verifies the About page renders with its title.
func TestAboutPageTitle(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "about/index.html")
	helpers.AssertText(t, doc, ".gh-title", "About")
}

// TestAboutPageContent verifies the About page body is rendered.
func TestAboutPageContent(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "about/index.html")
	helpers.AssertSelectorAtLeast(t, doc, ".gh-content", 1)
}

// ----------------------------------------------------------------------------
// M9 — Category term page
// ----------------------------------------------------------------------------

// TestCategoryTechPageExists verifies /categories/technology/ is built.
func TestCategoryTechPageExists(t *testing.T) {
	buildOnce(t)
	helpers.FileExists(t, "categories/technology/index.html")
}

// TestCategoryTechPageTitle verifies the term page shows the category name.
func TestCategoryTechPageTitle(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "categories/technology/index.html")
	helpers.AssertText(t, doc, ".gh-page-head h1", "Technology")
}

// TestCategoryTechPageHasPosts verifies the category page shows posts for that category.
func TestCategoryTechPageHasPosts(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "categories/technology/index.html")
	// "Technology" has 2 fixture posts
	helpers.AssertSelectorAtLeast(t, doc, ".gh-card", 1)
}

// TestCategoryTechDescription verifies the category description from _index.md is shown.
func TestCategoryTechDescription(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "categories/technology/index.html")
	helpers.AssertText(t, doc, ".gh-page-head p", "Computer technology")
}

// TestCategoryAIPageTitle verifies the AI category page.
func TestCategoryAIPageTitle(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "categories/ai/index.html")
	helpers.AssertText(t, doc, ".gh-page-head h1", "AI")
}

// ----------------------------------------------------------------------------
// M10 — 404 page
// ----------------------------------------------------------------------------

// Test404PageExists verifies 404.html is in the built output.
func Test404PageExists(t *testing.T) {
	buildOnce(t)
	helpers.FileExists(t, "404.html")
}

// Test404HasErrorCode verifies the 404 code is visible.
func Test404HasErrorCode(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "404.html")
	helpers.AssertText(t, doc, ".gh-error-code", "404")
}

// Test404HasHomeLink verifies there is a "Go back home" link.
func Test404HasHomeLink(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "404.html")
	helpers.AssertSelector(t, doc, ".gh-error-link[href]", 1)
}

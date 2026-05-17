package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// Taxonomy conventions
//
// The exampleSite declares:
//
//   [taxonomies]
//     category = "category"
//     tag      = "tag"
//
// so the front-matter keys are `category:` / `tag:` and the URL segments
// are /category/<slug>/ and /tag/<slug>/. These tests pin that contract
// across every layout that emits a taxonomy link, so an accidental
// regression to /categories/ or /tags/ fails immediately.
// ----------------------------------------------------------------------------

const (
	categoryURLPrefix = "/category/"
	tagURLPrefix      = "/tag/"
)

// TestTaxonomyTermPagesUseSingularURL verifies category and tag term
// pages are emitted at the singular URL.
func TestTaxonomyTermPagesUseSingularURL(t *testing.T) {
	buildOnce(t)
	helpers.FileExists(t, "category/index.html")
	helpers.FileExists(t, "category/technology/index.html")
	helpers.FileExists(t, "tag/index.html")
	helpers.FileExists(t, "tag/rust/index.html")
}

// TestTaxonomyPluralURLsAreNotBuilt guards against a regression where a
// plural-form taxonomy directory leaks back into the build (e.g. someone
// adds a stale `categories = "categories"` line).
func TestTaxonomyPluralURLsAreNotBuilt(t *testing.T) {
	buildOnce(t)

	for _, p := range []string{"categories/index.html", "tags/index.html"} {
		full := filepath.Join(helpers.PublicDir(), filepath.FromSlash(p))
		if _, err := os.Stat(full); err == nil {
			t.Errorf("unexpected plural taxonomy page built: %s", p)
		}
	}
}

// TestSinglePostPrimaryCategoryHref verifies the primary category link on
// a single post page points at the exact /category/<slug>/ URL emitted
// by Hugo for the rust post's "Technology" category.
func TestSinglePostPrimaryCategoryHref(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "posts/rust-ownership-model/index.html")

	href := helpers.Href(doc, ".cat-links a")
	const want = "/category/technology/"
	if href != want {
		t.Errorf("category href: got %q, want %q", href, want)
	}
}

// TestSinglePostTagHrefs verifies every tag link on a single post page
// points at /tag/<slug>/ (and not /tags/).
func TestSinglePostTagHrefs(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "posts/rust-ownership-model/index.html")

	hrefs := helpers.AllHrefs(doc, ".entry-footer .tags-links a")
	if len(hrefs) == 0 {
		t.Fatal("expected at least one tag link in entry footer")
	}

	for _, href := range hrefs {
		if !strings.HasPrefix(href, tagURLPrefix) {
			t.Errorf("tag href: got %q, want prefix %q", href, tagURLPrefix)
		}
		if strings.HasPrefix(href, "/tags/") {
			t.Errorf("tag href: got plural %q, want singular %q...", href, tagURLPrefix)
		}
	}
}

// TestSinglePostTagsAllRendered verifies *all* tags from front matter
// appear as links — guards against the .GetTerms call returning only one.
func TestSinglePostTagsAllRendered(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "posts/rust-ownership-model/index.html")

	// Rust post front matter: tag: ["rust", "programming", "memory"]
	want := []string{"rust", "programming", "memory"}

	hrefs := helpers.AllHrefs(doc, ".entry-footer .tags-links a")
	gotByTag := map[string]bool{}
	for _, h := range hrefs {
		gotByTag[h] = true
	}

	for _, slug := range want {
		expected := tagURLPrefix + slug + "/"
		if !gotByTag[expected] {
			t.Errorf("missing tag link %q (got: %v)", expected, hrefs)
		}
	}
}

// TestCardCategoryHrefIsSingular verifies the homepage cards' category
// links also use /category/.
func TestCardCategoryHrefIsSingular(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	hrefs := helpers.AllHrefs(doc, ".gh-card .gh-card-tag")
	if len(hrefs) == 0 {
		t.Fatal("expected at least one .gh-card-tag href on homepage")
	}

	for _, href := range hrefs {
		if !strings.HasPrefix(href, categoryURLPrefix) {
			t.Errorf("card category href: got %q, want prefix %q", href, categoryURLPrefix)
		}
	}
}

// TestCardCategoryHrefOnCategoryTermPage verifies cards rendered on a
// category term page (/category/technology/) also emit /category/...
// hrefs. Term pages reuse card.html, so this guards against a regression
// where the card's .Page context resolves taxonomy URLs differently
// outside the homepage.
func TestCardCategoryHrefOnCategoryTermPage(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "category/technology/index.html")

	hrefs := helpers.AllHrefs(doc, ".gh-card .gh-card-tag")
	if len(hrefs) == 0 {
		t.Fatal("expected at least one .gh-card-tag href on /category/technology/")
	}
	for _, href := range hrefs {
		if !strings.HasPrefix(href, categoryURLPrefix) {
			t.Errorf("term-page card category href: got %q, want prefix %q",
				href, categoryURLPrefix)
		}
	}
}

// TestCardCategoryHrefOnTagTermPage verifies cards rendered on a tag
// term page (/tag/rust/) emit /category/... hrefs (the card always
// shows the post's primary *category*, regardless of which term page
// the listing belongs to).
func TestCardCategoryHrefOnTagTermPage(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "tag/rust/index.html")

	hrefs := helpers.AllHrefs(doc, ".gh-card .gh-card-tag")
	if len(hrefs) == 0 {
		t.Fatal("expected at least one .gh-card-tag href on /tag/rust/")
	}
	for _, href := range hrefs {
		if !strings.HasPrefix(href, categoryURLPrefix) {
			t.Errorf("tag-term-page card category href: got %q, want prefix %q",
				href, categoryURLPrefix)
		}
	}
}

// TestSidebarMenuHrefsAreSingular pins the hugo.toml [menu] URLs to the
// singular form. (The sidebar test file already checks the exact list,
// but this is a higher-level invariant that fails fast when someone
// reverts the menu config.)
func TestSidebarMenuHrefsAreSingular(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	hrefs := helpers.AllHrefs(doc, "#category-menu-tab .cat-item a")
	if len(hrefs) == 0 {
		t.Fatal("expected sidebar category menu items")
	}
	for _, href := range hrefs {
		if !strings.HasPrefix(href, categoryURLPrefix) {
			t.Errorf("sidebar menu href: got %q, want prefix %q", href, categoryURLPrefix)
		}
	}
}

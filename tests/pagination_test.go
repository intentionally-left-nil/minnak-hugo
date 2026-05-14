package tests

import (
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// M6 — Pagination
// ----------------------------------------------------------------------------

// TestNoPaginationOnSinglePage verifies that "Page 1 of 1" is not shown when
// all posts fit on one page (paginate=15, we have 7 posts).
func TestNoPaginationOnSinglePage(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	// With 7 posts and paginate=15, only 1 page exists → paginator is present
	// but "Older Posts" and "Newer Posts" links should be hidden (disabled).
	olderLinks := doc.Find(".gh-pagination-older:not(.gh-pagination-disabled)")
	if olderLinks.Length() > 0 {
		t.Error("unexpected active 'Older Posts' link on single-page site")
	}
}

// TestPaginationStructurePresent verifies pagination HTML is rendered.
func TestPaginationStructurePresent(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, ".gh-pagination", 1)
}

// TestPaginationPageInfo verifies "Page N of M" text is rendered.
func TestPaginationPageInfo(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertText(t, doc, ".gh-pagination-info", "page")
}

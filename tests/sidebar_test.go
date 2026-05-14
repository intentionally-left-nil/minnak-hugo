package tests

import (
	"strings"
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// M2.1 — Sidebar markup parity
// ----------------------------------------------------------------------------

// TestSidebarExists verifies the sidebar container is present.
func TestSidebarExists(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, "#left-sidebar", 1)
}

// TestSidebarHasTwoTabs verifies exactly two icon tabs are rendered.
func TestSidebarHasTwoTabs(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, "#left-sidebar .vertical-menu-item", 2)
}

// TestSidebarFirstTabIsActive verifies the categories tab is active by default.
func TestSidebarFirstTabIsActive(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, "#left-sidebar .vertical-menu-item.active-menu", 1)

	// The active item must be the first one (category-nav).
	first := doc.Find("#left-sidebar .vertical-menu-item").First()
	if !first.HasClass("active-menu") {
		t.Error("expected the first .vertical-menu-item to have active-menu class")
	}
}

// TestSidebarTabIDs verifies tab anchor and panel IDs match the spec.
func TestSidebarTabIDs(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	helpers.AssertSelector(t, doc, "#v-pills-cats-tab", 1)
	helpers.AssertSelector(t, doc, "#category-menu-tab", 1)
	helpers.AssertSelector(t, doc, "#v-pills-widgets-tab", 1)
	helpers.AssertSelector(t, doc, "#widget-tab", 1)
}

// TestSidebarToggleButton verifies the mobile toggle button is present.
func TestSidebarToggleButton(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, "#toggle-button.toggle-button", 1)
}

// TestSidebarFooterExists verifies the empty sidebar footer placeholder is rendered.
func TestSidebarFooterExists(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, "#left-sidebar-footer", 1)
}

// ----------------------------------------------------------------------------
// M2.2 — Categories tab driven by Site.Menus.main
// ----------------------------------------------------------------------------

var wantMenuItems = []struct {
	name string
	href string
}{
	{"Technology", "/categories/technology/"},
	{"AI", "/categories/ai/"},
	{"Cantonese", "/categories/cantonese/"},
	{"Photography", "/categories/photography/"},
	{"PNW", "/categories/pnw/"},
}

// TestSidebarCategoryMenuItems verifies the sidebar categories list
// is populated from Site.Menus.main.
func TestSidebarCategoryMenuItems(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	items := doc.Find("#category-menu-tab .cat-item a")
	if items.Length() != len(wantMenuItems) {
		t.Errorf("sidebar category items: got %d, want %d", items.Length(), len(wantMenuItems))
	}

	for i, want := range wantMenuItems {
		item := items.Eq(i)
		text := strings.TrimSpace(item.Text())
		href, _ := item.Attr("href")

		if text != want.name {
			t.Errorf("menu item %d: name = %q, want %q", i, text, want.name)
		}
		if href != want.href {
			t.Errorf("menu item %d: href = %q, want %q", i, href, want.href)
		}
	}
}

// ----------------------------------------------------------------------------
// M2.3 — Recent posts list (top 5 by date, Type=posts)
// ----------------------------------------------------------------------------

// TestSidebarRecentPostsCount verifies exactly 5 recent posts are listed.
func TestSidebarRecentPostsCount(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	items := doc.Find(".widget_recent_entries li")
	if items.Length() != 5 {
		t.Errorf("recent posts count: got %d, want 5", items.Length())
	}
}

// TestSidebarRecentPostsAreLinks verifies each recent post is a link.
func TestSidebarRecentPostsAreLinks(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, ".widget_recent_entries li a[href]", 5)
}

// TestSidebarRecentPostsOrder verifies the most-recent post is listed first.
// The fixture most-recent post (by date) is "Rust's Ownership Model..." (2026-03-15).
func TestSidebarRecentPostsOrder(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	firstTitle := strings.TrimSpace(
		doc.Find(".widget_recent_entries li a").First().Text(),
	)
	wantFirst := "Rust's Ownership Model Is Actually About Time"
	if firstTitle != wantFirst {
		t.Errorf("first recent post: got %q, want %q", firstTitle, wantFirst)
	}
}

// ----------------------------------------------------------------------------
// M2.4 — Pagefind mount point & data-pagefind-body on single pages
// ----------------------------------------------------------------------------

// TestSidebarSearchMountExists verifies the #search mount point is in the widget tab.
func TestSidebarSearchMountExists(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, "#widget-tab #search", 1)
}

// TestSinglePostPagefindBody verifies data-pagefind-body is on the article.
func TestSinglePostPagefindBody(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "posts/rust-ownership-model/index.html")
	helpers.AssertSelector(t, doc, "article[data-pagefind-body]", 1)
}

// TestSinglePostPagefindMeta verifies data-pagefind-meta="title" is present.
func TestSinglePostPagefindMeta(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "posts/rust-ownership-model/index.html")
	helpers.AssertSelector(t, doc, "[data-pagefind-meta='title']", 1)
}

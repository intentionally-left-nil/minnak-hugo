package tests

import (
	"strings"
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// Card cover image
//
// card.html resolves the post's hero image through a single mechanism:
//   1. The resource path stored in cover.src front matter.
//   2. No image at all — the <picture> block is skipped entirely.
//
// cover.alt provides the alt text; if absent it falls back to the page title.
//
// These tests pin both paths from the homepage card listing so a
// regression in either branch fails immediately.
// ----------------------------------------------------------------------------

// findCardForPost returns the .gh-card <article> that links to the given
// post permalink, or fails the test if none is found.
//
// The homepage renders three copies of every card (3-col, 2-col, 1-col
// grids) plus a sidebar "recent posts" list that *also* mentions every
// post URL. To get a usable card chunk we scan only between
// `<article class="gh-card">` and `</article>`.
func findCardForPost(t *testing.T, indexHTML, postSlug string) string {
	t.Helper()

	const open = `<article class="gh-card">`
	want := "/posts/" + postSlug + "/"

	rest := indexHTML
	for {
		i := strings.Index(rest, open)
		if i < 0 {
			break
		}
		rest = rest[i:]
		j := strings.Index(rest, "</article>")
		if j < 0 {
			break
		}
		card := rest[:j+len("</article>")]
		if strings.Contains(card, want) {
			return card
		}
		rest = rest[j+len("</article>"):]
	}
	t.Fatalf("no homepage card linking to %q", want)
	return ""
}

// readIndexHTML reads the homepage as a string (for substring scans).
func readIndexHTML(t *testing.T) string {
	t.Helper()
	doc := helpers.ParseFile(t, "index.html")
	html, err := doc.Html()
	if err != nil {
		t.Fatalf("could not serialize index.html: %v", err)
	}
	return html
}

// TestCardCoverSrcWithBundleFile verifies cards for posts with a cover.src
// pointing at a page-bundle file render the <picture> block sourcing from
// /posts/<slug>/ (Hugo's processed-image output).
func TestCardCoverSrcWithBundleFile(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	// rust-ownership-model has cover.src: "feature.jpg" in its bundle.
	card := findCardForPost(t, html, "rust-ownership-model")

	if !strings.Contains(card, "<picture") {
		t.Error("expected <picture> in card for rust post (cover.src: feature.jpg)")
	}
	if !strings.Contains(card, `srcset="/posts/rust-ownership-model/`) {
		t.Errorf("expected srcset to source from /posts/rust-ownership-model/, card was:\n%s", card)
	}
	if !strings.Contains(card, "image/webp") {
		t.Error("expected webp <source> in card")
	}
}

// TestCardCoverSrcWithExplicitPath verifies a post with cover.src pointing
// at an explicit path (not "feature.jpg") also renders a <picture>.
func TestCardCoverSrcWithExplicitPath(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	// fairphone-review has cover.src: "images/Fairphone.jpg"
	card := findCardForPost(t, html, "fairphone-review")

	if !strings.Contains(card, "<picture") {
		t.Errorf("expected <picture> in fairphone card (cover.src: images/Fairphone.jpg). Card was:\n%s", card)
	}
	// The processed image filename starts with the original basename.
	if !strings.Contains(card, "/Fairphone") {
		t.Errorf("expected processed Fairphone image in srcset/src. Card was:\n%s", card)
	}
}

// TestCardCoverAltUsed verifies the alt text on a card comes from cover.alt.
func TestCardCoverAltUsed(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	card := findCardForPost(t, html, "fairphone-review")
	wantAlt := "A Fairphone on a wooden desk with its back cover removed"
	if !strings.Contains(card, `alt="`+wantAlt+`"`) {
		t.Errorf("expected alt=%q on fairphone card. Card was:\n%s", wantAlt, card)
	}
}

// TestCardCoverAltUsedForBundleImage verifies that cover.alt is used even
// when the image is a page-bundle file (not an explicit images/ path).
func TestCardCoverAltUsedForBundleImage(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	// rust-ownership-model has cover.src: "feature.jpg" and
	// cover.alt: "Photo by Naoki Suzuki on Unsplash".
	card := findCardForPost(t, html, "rust-ownership-model")
	wantAlt := "Photo by Naoki Suzuki on Unsplash"
	if !strings.Contains(card, `alt="`+wantAlt+`"`) {
		t.Errorf("expected alt=%q on rust card. Card was:\n%s", wantAlt, card)
	}
}

// TestCardNoImageRendersGracefully verifies cards with no cover.src still
// render their text content (title, summary, read-more) but skip the
// <picture> block.
func TestCardNoImageRendersGracefully(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	// gallery-mount-rainier has no cover.src.
	card := findCardForPost(t, html, "gallery-mount-rainier")

	if strings.Contains(card, "<picture") {
		t.Errorf("did not expect <picture> in card with no cover.src. Card was:\n%s", card)
	}
	if !strings.Contains(card, "Mount Rainier in Four Frames") {
		t.Error("expected gallery card to still show its title")
	}
	if !strings.Contains(card, "read-more-link") {
		t.Error("expected gallery card to still show its read-more link")
	}
}

// TestCardCoverImageProcessesToWebP verifies the cover.src pipeline emits
// WebP+JPEG <source> elements with a six-width srcset.
func TestCardCoverImageProcessesToWebP(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	card := findCardForPost(t, html, "fairphone-review")

	if !strings.Contains(card, `type="image/webp"`) {
		t.Errorf("expected image/webp <source> on cover.src card. Card was:\n%s", card)
	}
	if !strings.Contains(card, "30w") || !strings.Contains(card, "2000w") {
		t.Errorf("expected six-width srcset on cover.src card. Card was:\n%s", card)
	}
}

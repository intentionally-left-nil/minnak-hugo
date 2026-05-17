package tests

import (
	"strings"
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// Card cover image fallback
//
// card.html resolves the post's hero image through three branches:
//   1. A page-bundle resource matching `feature.*` (the canonical case,
//      used by most fixture posts).
//   2. The path stored in `cover.image` front matter (the WordPress
//      migration case — see fixture posts/fairphone-review/).
//   3. No image at all — the <picture> block is skipped entirely.
//
// These tests pin all three paths from the homepage card listing so a
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

// TestCardFeatureImagePath verifies cards for posts with a feature.*
// resource render the <picture> block sourcing from /posts/<slug>/
// (Hugo's processed-image output).
func TestCardFeatureImagePath(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	// rust-ownership-model has feature.jpg in its bundle.
	card := findCardForPost(t, html, "rust-ownership-model")

	if !strings.Contains(card, "<picture") {
		t.Error("expected <picture> in card for rust post (feature.jpg present)")
	}
	if !strings.Contains(card, `srcset="/posts/rust-ownership-model/`) {
		t.Errorf("expected srcset to source from /posts/rust-ownership-model/, card was:\n%s", card)
	}
	if !strings.Contains(card, "image/webp") {
		t.Error("expected webp <source> in card")
	}
}

// TestCardCoverImageFallback verifies a post with no feature.* but a
// `cover.image` path still gets a rendered <picture>.
func TestCardCoverImageFallback(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	// fairphone-review has cover.image: images/Fairphone.jpg, no feature.*
	card := findCardForPost(t, html, "fairphone-review")

	if !strings.Contains(card, "<picture") {
		t.Errorf("expected <picture> in fairphone card (cover.image fallback). Card was:\n%s", card)
	}
	// The processed image filename starts with the original basename.
	if !strings.Contains(card, "/Fairphone") {
		t.Errorf("expected processed Fairphone image in srcset/src. Card was:\n%s", card)
	}
}

// TestCardCoverAltUsed verifies the alt text on the cover-fallback path
// comes from cover.alt (not the page title).
func TestCardCoverAltUsed(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	card := findCardForPost(t, html, "fairphone-review")
	wantAlt := "A Fairphone on a wooden desk with its back cover removed"
	if !strings.Contains(card, `alt="`+wantAlt+`"`) {
		t.Errorf("expected alt=%q on fairphone card. Card was:\n%s", wantAlt, card)
	}
}

// TestCardAltFallsBackToTitle verifies that when neither cover.alt nor
// (the deprecated) feature_image_alt is set, the alt text falls back to
// the post title.
func TestCardAltFallsBackToTitle(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	// rust-ownership-model has no cover.alt — alt should be the post title.
	card := findCardForPost(t, html, "rust-ownership-model")

	wantTitle := "Rust&#39;s Ownership Model Is Actually About Time"
	if !strings.Contains(card, `alt="`+wantTitle+`"`) {
		t.Errorf("expected alt to fall back to post title. Card was:\n%s", card)
	}
}

// TestCardNoImageRendersGracefully verifies cards with neither feature.*
// nor cover.image still render their text content (title, summary,
// read-more) but skip the <picture> block.
func TestCardNoImageRendersGracefully(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	// gallery-mount-rainier has no feature.* and no cover.image.
	card := findCardForPost(t, html, "gallery-mount-rainier")

	if strings.Contains(card, "<picture") {
		t.Errorf("did not expect <picture> in card with no image. Card was:\n%s", card)
	}
	if !strings.Contains(card, "Mount Rainier in Four Frames") {
		t.Error("expected gallery card to still show its title")
	}
	if !strings.Contains(card, "read-more-link") {
		t.Error("expected gallery card to still show its read-more link")
	}
}

// TestCardCoverImageProcessesToWebP verifies the cover.image fallback
// path runs through the same WebP+JPEG <source> pipeline as feature.*.
func TestCardCoverImageProcessesToWebP(t *testing.T) {
	buildOnce(t)
	html := readIndexHTML(t)

	card := findCardForPost(t, html, "fairphone-review")

	if !strings.Contains(card, `type="image/webp"`) {
		t.Errorf("expected image/webp <source> on cover.image fallback. Card was:\n%s", card)
	}
	if !strings.Contains(card, "30w") || !strings.Contains(card, "2000w") {
		t.Errorf("expected six-width srcset on cover.image fallback. Card was:\n%s", card)
	}
}

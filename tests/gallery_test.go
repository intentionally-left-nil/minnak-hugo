package tests

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// Gallery shortcode markup
//
// The shortcode {{< gallery cols="3" >}} wraps nested {{< figure >}}
// shortcodes. Each child figure detects .Parent and renders as a
// <figure class="gallery-item"> with a 400×400 square thumbnail (cropped
// via .Fill) linked to the full-size image. alt and caption are separate
// params on each figure — they can hold different values.
//
// These tests run against the `posts/gallery-mount-rainier/` fixture.
// Responsive layout (column counts at different viewports) lives in the
// Playwright suite — these are pure markup assertions.
// ----------------------------------------------------------------------------

const galleryPostPath = "posts/gallery-mount-rainier/index.html"

// TestGalleryWrapperColsClass verifies the wrapper carries a
// gallery-cols-N class derived from the cols= parameter.
func TestGalleryWrapperColsClass(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	// The fixture invokes {{< gallery cols="3" >}}.
	helpers.AssertSelector(t, doc, ".gallery.gallery-cols-3", 1)
}

// TestGalleryRendersAllImages verifies one <figure> per nested figure
// shortcode (the fixture has four).
func TestGalleryRendersAllImages(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	helpers.AssertSelector(t, doc, ".gallery .gallery-item", 4)
	helpers.AssertSelector(t, doc, ".gallery .gallery-item img", 4)
	helpers.AssertSelector(t, doc, ".gallery .gallery-item a[href]", 4)
}

// TestGalleryThumbnailsLinkToFullImage verifies each thumbnail wraps in
// an <a> that points at the original (unprocessed) image, while the
// <img src> uses the processed thumbnail.
func TestGalleryThumbnailsLinkToFullImage(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	hrefs := helpers.AllHrefs(doc, ".gallery .gallery-item a")
	if len(hrefs) != 4 {
		t.Fatalf("expected 4 thumbnail links, got %d", len(hrefs))
	}

	// Every link should point at the page-bundle image path. The fixture
	// images are named photo-01.jpg through photo-04.jpg.
	for i, href := range hrefs {
		want := "/posts/gallery-mount-rainier/images/gallery/photo-0"
		if !strings.HasPrefix(href, want) {
			t.Errorf("thumbnail link %d: got %q, want prefix %q", i, href, want)
		}
		if !strings.HasSuffix(href, ".jpg") {
			t.Errorf("thumbnail link %d: got %q, want .jpg suffix", i, href)
		}
	}

	// And every <img src> should reference Hugo's processed-image output
	// (filenames contain a hash signature like "_hu_<hash>").
	imgs := doc.Find(".gallery .gallery-item img")
	imgs.Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if !strings.Contains(src, "_hu_") {
			t.Errorf("thumbnail %d src=%q: expected processed-image filename (contains _hu_)", i, src)
		}
	})
}

// TestGalleryPreservesInlineOrder verifies the shortcode emits figures in
// the order they appear inline in the content (photo-01 through photo-04).
func TestGalleryPreservesInlineOrder(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	hrefs := helpers.AllHrefs(doc, ".gallery .gallery-item a")
	wantOrder := []string{"photo-01", "photo-02", "photo-03", "photo-04"}

	if len(hrefs) != len(wantOrder) {
		t.Fatalf("expected %d gallery links, got %d", len(wantOrder), len(hrefs))
	}
	for i, want := range wantOrder {
		if !strings.Contains(hrefs[i], want) {
			t.Errorf("gallery position %d: href=%q, want to contain %q", i, hrefs[i], want)
		}
	}
}

// TestGalleryCaptionsFromShortcodeParam verifies <figcaption> text is
// pulled from the caption= param on each nested figure shortcode.
func TestGalleryCaptionsFromShortcodeParam(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	captions := helpers.AllTexts(doc, ".gallery .gallery-item figcaption")
	wantCaptions := []string{
		"Sunrise on Sunrise",
		"Tipsoo Lake at dawn",
		"Skyline Trail wildflowers",
		"Nisqually glacier from Glacier Vista",
	}

	if len(captions) != len(wantCaptions) {
		t.Fatalf("expected %d captions, got %d: %v", len(wantCaptions), len(captions), captions)
	}
	for i, want := range wantCaptions {
		if captions[i] != want {
			t.Errorf("caption %d: got %q, want %q", i, captions[i], want)
		}
	}
}

// TestGalleryThumbnailsAreSquare verifies the .Fill processor produced
// square thumbnails (the requested 400x400 crop).
func TestGalleryThumbnailsAreSquare(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	imgs := doc.Find(".gallery .gallery-item img")
	if imgs.Length() == 0 {
		t.Fatal("no gallery images found")
	}

	imgs.Each(func(i int, s *goquery.Selection) {
		w, _ := s.Attr("width")
		h, _ := s.Attr("height")
		if w == "" || h == "" {
			t.Errorf("img %d: missing width/height attributes", i)
			return
		}
		if w != h {
			t.Errorf("img %d: got %sx%s, expected square thumbnail", i, w, h)
		}
		if w != "400" {
			t.Errorf("img %d: got width=%s, expected 400", i, w)
		}
	})
}

// TestGalleryImgAltIsNonEmpty verifies each gallery thumbnail has a
// non-empty alt attribute.
func TestGalleryImgAltIsNonEmpty(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	imgs := doc.Find(".gallery .gallery-item img")
	imgs.Each(func(i int, s *goquery.Selection) {
		alt, _ := s.Attr("alt")
		if alt == "" {
			t.Errorf("img %d: empty alt attribute", i)
		}
	})
}

// TestGalleryAltAndCaptionAreDistinct verifies that the alt attribute and
// the <figcaption> text are independent fields and can hold different values.
// The fixture uses shorter display captions with fuller alt descriptions.
func TestGalleryAltAndCaptionAreDistinct(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	items := doc.Find(".gallery .gallery-item")
	items.Each(func(i int, s *goquery.Selection) {
		alt, _ := s.Find("img").First().Attr("alt")
		caption := strings.TrimSpace(s.Find("figcaption").Text())
		if alt == "" || caption == "" {
			t.Errorf("item %d: alt=%q caption=%q — both must be non-empty", i, alt, caption)
			return
		}
		if alt == caption {
			t.Errorf("item %d: alt and caption are identical (%q) — fixture expects them to differ", i, alt)
		}
	})
}

// TestGalleryAltTexts verifies the exact alt text on each gallery image
// matches the alt= param passed to each nested figure shortcode.
func TestGalleryAltTexts(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	imgs := doc.Find(".gallery .gallery-item img")
	wantAlts := []string{
		"Sunrise on east-facing slopes of Mount Rainier catching the first light",
		"Tipsoo Lake reflecting Mount Rainier at dawn",
		"Lupine and paintbrush wildflowers on the Skyline Trail",
		"The Nisqually glacier viewed from Glacier Vista",
	}

	if imgs.Length() != len(wantAlts) {
		t.Fatalf("expected %d images, got %d", len(wantAlts), imgs.Length())
	}
	imgs.Each(func(i int, s *goquery.Selection) {
		alt, _ := s.Attr("alt")
		if alt != wantAlts[i] {
			t.Errorf("img %d alt: got %q, want %q", i, alt, wantAlts[i])
		}
	})
}

// TestGalleryLazyLoading verifies thumbnails opt into lazy loading.
func TestGalleryLazyLoading(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	imgs := doc.Find(".gallery .gallery-item img")
	imgs.Each(func(i int, s *goquery.Selection) {
		loading, _ := s.Attr("loading")
		if loading != "lazy" {
			t.Errorf("img %d: loading=%q, want \"lazy\"", i, loading)
		}
	})
}

// TestGalleryShortcodeNotPresentOnNonGalleryPost verifies the gallery-related
// selectors are not matched on a post that does not use the gallery shortcode.
func TestGalleryShortcodeNotPresentOnNonGalleryPost(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "posts/rust-ownership-model/index.html")

	if doc.Find(".gallery").Length() != 0 {
		t.Error("did not expect a .gallery element on a non-gallery post")
	}
}

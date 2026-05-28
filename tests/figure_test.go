package tests

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// Figure shortcode (layouts/shortcodes/figure.html)
//
// The shortcode has three modes:
//
//  1. cover=true — reads src/alt/caption from the page's cover: front matter.
//     Explicitly passed params override front-matter values.
//     Renders a standalone <figure> with a full <picture> srcset pipeline.
//
//  2. Standalone — src/alt/caption passed directly as shortcode params.
//     Renders a standalone <figure> with a full <picture> srcset pipeline.
//
//  3. Gallery child — called inside {{< gallery >}}.
//     Detects .Parent and renders <figure class="gallery-item"> with a
//     400×400 square-cropped thumbnail linked to the full-size original.
//
// Fixtures:
//
//   - posts/carbon-river-golden-hour uses {{< figure cover=true >}} (no
//     caption) — exercises mode 1 and the full image-processing pipeline.
//   - posts/gallery-mount-rainier uses nested {{< figure >}} inside
//     {{< gallery >}} — exercises mode 3. alt ≠ caption in the fixture.
//
// Gallery-specific tests (thumbnail dimensions, ordering, etc.) live in
// gallery_test.go. The tests here focus on the figure shortcode itself.
// ----------------------------------------------------------------------------

const carbonRiverPath = "posts/carbon-river-golden-hour/index.html"

// ── Standalone / cover=true mode ─────────────────────────────────────────────

// TestFigureRendersElement verifies {{< figure cover=true >}} emits a
// <figure> element in the post body.
func TestFigureRendersElement(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)
	helpers.AssertSelectorAtLeast(t, doc, ".entry-content figure", 1)
}

// TestFigureRendersPicture verifies the standalone figure contains a
// <picture> element (not a bare <img>).
func TestFigureRendersPicture(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)
	helpers.AssertSelectorAtLeast(t, doc, ".entry-content figure picture", 1)
}

// TestFigureWebPSource verifies the figure renders a
// <source type="image/webp"> for modern browsers.
func TestFigureWebPSource(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)
	helpers.AssertSelectorAtLeast(t, doc, `.entry-content figure source[type="image/webp"]`, 1)
}

// TestFigureJPEGFallbackSource verifies the figure renders a second <source>
// without a type attribute (the JPEG fallback for older browsers).
func TestFigureJPEGFallbackSource(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)

	// There should be exactly two <source> elements: one WebP, one JPEG.
	helpers.AssertSelector(t, doc, ".entry-content figure source", 2)

	// The JPEG source has no type attribute.
	jpegSrc := doc.Find(".entry-content figure source:not([type])")
	if jpegSrc.Length() == 0 {
		t.Error("expected a <source> without a type attribute (JPEG fallback)")
	}
}

// TestFigureSrcsetWidths verifies both sources carry the six expected width
// descriptors (480w through 2000w).
func TestFigureSrcsetWidths(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)

	wantWidths := []string{"480w", "720w", "960w", "1200w", "1600w", "2000w"}
	doc.Find(".entry-content figure source").Each(func(i int, s *goquery.Selection) {
		srcset, _ := s.Attr("srcset")
		for _, w := range wantWidths {
			if !strings.Contains(srcset, w) {
				t.Errorf("source %d srcset missing %q: %s", i, w, srcset)
			}
		}
	})
}

// TestFigureImgLazyLoading verifies loading="lazy" is set on the <img>
// fallback element.
func TestFigureImgLazyLoading(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)
	helpers.AssertAttr(t, doc, ".entry-content figure img", "loading", "lazy")
}

// TestFigureImgHasDimensions verifies explicit width and height attributes
// are present on the <img> to allow the browser to reserve layout space
// before the image loads (prevents CLS).
func TestFigureImgHasDimensions(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)

	img := doc.Find(".entry-content figure img").First()
	if img.Length() == 0 {
		t.Fatal("no <img> found inside figure")
	}
	if w, _ := img.Attr("width"); w == "" {
		t.Error("figure <img> missing width attribute")
	}
	if h, _ := img.Attr("height"); h == "" {
		t.Error("figure <img> missing height attribute")
	}
}

// TestFigureCoverAltUsed verifies that cover=true reads the alt text from
// the page's cover.alt front matter.
func TestFigureCoverAltUsed(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)

	// carbon-river-golden-hour has cover.alt: "Photo by Satwika Ananta on Unsplash"
	helpers.AssertAttr(t, doc, ".entry-content figure img", "alt", "Photo by Satwika Ananta on Unsplash")
}

// TestFigureNoCaptionByDefault verifies no <figcaption> is rendered when
// neither the shortcode param nor cover.caption is set.
// carbon-river-golden-hour has no cover.caption.
func TestFigureNoCaptionByDefault(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)

	if doc.Find(".entry-content figure figcaption").Length() != 0 {
		t.Error("expected no <figcaption> when no caption is set")
	}
}

// TestFigureNotGalleryItem verifies a standalone figure does not carry the
// gallery-item class (that class is only for gallery child mode).
func TestFigureNotGalleryItem(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)

	if doc.Find(".entry-content figure.gallery-item").Length() != 0 {
		t.Error("standalone figure should not have gallery-item class")
	}
}

// TestFigureSrcFromCoverFrontMatter verifies the figure srcset references
// the file named in cover.src (feature.jpg) from the page bundle.
func TestFigureSrcFromCoverFrontMatter(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)

	// Processed images keep "feature" in the filename.
	srcset, _ := doc.Find(".entry-content figure source[type='image/webp']").First().Attr("srcset")
	if !strings.Contains(srcset, "/posts/carbon-river-golden-hour/feature") {
		t.Errorf("expected srcset to reference carbon-river feature image. Got: %s", srcset)
	}
}

// ── Gallery child mode ────────────────────────────────────────────────────────

// TestFigureGalleryItemClass verifies that figure shortcodes nested inside
// {{< gallery >}} carry the gallery-item class.
func TestFigureGalleryItemClass(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)
	helpers.AssertSelector(t, doc, ".gallery figure.gallery-item", 4)
}

// TestFigureGalleryItemHasNoSrcset verifies that gallery thumbnails render
// as plain <img> elements (not <picture> with srcset) — they use .Fill for
// a fixed square crop rather than the responsive srcset pipeline.
func TestFigureGalleryItemHasNoSrcset(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)

	if doc.Find(".gallery figure.gallery-item picture").Length() != 0 {
		t.Error("gallery thumbnails should not contain a <picture> element")
	}
	if doc.Find(".gallery figure.gallery-item source").Length() != 0 {
		t.Error("gallery thumbnails should not contain <source> elements")
	}
}

// TestFigureGalleryItemHasCaption verifies that gallery figures render
// their caption= param as a <figcaption>.
func TestFigureGalleryItemHasCaption(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, galleryPostPath)
	helpers.AssertSelector(t, doc, ".gallery figure.gallery-item figcaption", 4)
}

// ── maxheight param ───────────────────────────────────────────────────────────

const maxheightPostPath = "posts/figure-maxheight-fixture/index.html"

// TestFigureMaxHeightStyleAttr verifies that a figure rendered with
// maxheight="400" emits style="max-height: 400px; width: auto;" on the <img>.
func TestFigureMaxHeightStyleAttr(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, maxheightPostPath)

	img := doc.Find(".entry-content figure img").First()
	if img.Length() == 0 {
		t.Fatal("no <img> found inside figure")
	}
	helpers.AssertStyleProp(t, doc, ".entry-content figure img", "max-height", "400px")
	helpers.AssertStyleProp(t, doc, ".entry-content figure img", "width", "auto")
}

// TestFigureMaxHeightPreservesDimensions verifies that width and height
// attributes are still present alongside the max-height style (so the browser
// can still reserve layout space before the image loads, preventing CLS).
func TestFigureMaxHeightPreservesDimensions(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, maxheightPostPath)

	img := doc.Find(".entry-content figure img").First()
	if img.Length() == 0 {
		t.Fatal("no <img> found inside figure")
	}
	if w, _ := img.Attr("width"); w == "" {
		t.Error("figure <img> with maxheight is missing the width attribute (needed for CLS prevention)")
	}
	if h, _ := img.Attr("height"); h == "" {
		t.Error("figure <img> with maxheight is missing the height attribute (needed for CLS prevention)")
	}
}

// TestFigureWithoutMaxHeightHasNoStyle verifies that a figure rendered
// without maxheight does NOT emit a style attribute on its <img>, so the
// feature is strictly opt-in.
func TestFigureWithoutMaxHeightHasNoStyle(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, carbonRiverPath)

	img := doc.Find(".entry-content figure img").First()
	if img.Length() == 0 {
		t.Fatal("no <img> found inside figure")
	}
	if style, exists := img.Attr("style"); exists && style != "" {
		t.Errorf("figure <img> without maxheight should have no style attribute, got %q", style)
	}
}

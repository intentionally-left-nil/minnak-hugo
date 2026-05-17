package tests

import (
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// ----------------------------------------------------------------------------
// M5 — Post grid (3/2/1-col responsive layout)
// ----------------------------------------------------------------------------

// TestGridHasThreeVariants verifies all three postfeed variants are in the DOM.
func TestGridHasThreeVariants(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	helpers.AssertSelector(t, doc, ".gh-postfeed-3col", 1)
	helpers.AssertSelector(t, doc, ".gh-postfeed-2col", 1)
	helpers.AssertSelector(t, doc, ".gh-postfeed-1col", 1)
}

// TestThreeColHasThreeColumns verifies 3-col layout has exactly 3 .gh-column divs.
func TestThreeColHasThreeColumns(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, ".gh-postfeed-3col .gh-column", 3)
}

// TestTwoColHasTwoColumns verifies 2-col layout has exactly 2 .gh-column divs.
func TestTwoColHasTwoColumns(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	helpers.AssertSelector(t, doc, ".gh-postfeed-2col .gh-column", 2)
}

// TestGridColumnDistribution verifies posts are distributed correctly across
// 3 columns using modular arithmetic: col i gets posts where index%3 == i.
func TestGridColumnDistribution(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	threeCol := doc.Find(".gh-postfeed-3col")
	cols := threeCol.Find(".gh-column")

	totalPosts := len(helpers.PostTitles())

	// Compute expected per-column counts from the modular distribution.
	wantCounts := []int{0, 0, 0}
	for i := 0; i < totalPosts; i++ {
		wantCounts[i%3]++
	}

	for i, want := range wantCounts {
		col := cols.Eq(i)
		got := col.Find(".gh-card").Length()
		if got != want {
			t.Errorf("3-col column %d: got %d cards, want %d (total posts=%d)",
				i, got, want, totalPosts)
		}
	}
}

// TestOneColHasAllPosts verifies the 1-col layout contains all posts.
func TestOneColHasAllPosts(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	// 1-col layout has cards directly (no gh-column wrapper)
	helpers.AssertSelector(t, doc, ".gh-postfeed-1col .gh-card", len(helpers.PostTitles()))
}

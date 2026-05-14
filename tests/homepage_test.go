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

// TestGridColumnDistribution verifies 7 posts are distributed correctly across
// 3 columns using modular arithmetic: col0=[0,3,6], col1=[1,4], col2=[2,5].
// (Post index 6 is in col 0 since 6%3==0; there is no post index 7.)
func TestGridColumnDistribution(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")

	threeCol := doc.Find(".gh-postfeed-3col")
	cols := threeCol.Find(".gh-column")

	// With 7 posts: col0 gets posts 0,3,6 → 3 cards
	//               col1 gets posts 1,4   → 2 cards
	//               col2 gets posts 2,5   → 2 cards
	wantCounts := []int{3, 2, 2}
	for i, want := range wantCounts {
		col := cols.Eq(i)
		got := col.Find(".gh-card").Length()
		if got != want {
			t.Errorf("3-col column %d: got %d cards, want %d", i, got, want)
		}
	}
}

// TestOneColHasAllPosts verifies the 1-col layout contains all 7 posts.
func TestOneColHasAllPosts(t *testing.T) {
	buildOnce(t)
	doc := helpers.ParseFile(t, "index.html")
	// 1-col layout has cards directly (no gh-column wrapper)
	helpers.AssertSelector(t, doc, ".gh-postfeed-1col .gh-card", len(helpers.PostTitles()))
}

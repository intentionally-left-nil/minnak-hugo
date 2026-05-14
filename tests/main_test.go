// TestMain is the test suite entry point. It runs a Hugo build once before
// all tests, so each test file can call helpers.Build(t, false, false) and
// reuse the already-built public/ directory.
package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// MINNAK_SKIP_BUILD=1 → skip the hugo build (use existing public/).
	// MINNAK_RUN_PAGEFIND=1 → also run npx pagefind after hugo.
	// Both flags are set by the CI workflow; omit them for fast local iteration.
	_ = os.Getenv("MINNAK_SKIP_BUILD")
	_ = os.Getenv("MINNAK_RUN_PAGEFIND")
	os.Exit(m.Run())
}

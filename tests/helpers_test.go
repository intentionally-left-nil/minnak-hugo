package tests

import (
	"testing"

	"github.com/intentionally-left-nil/minnak-hugo/tests/helpers"
)

// buildOnce is called at the top of every test that needs the built site.
// The helpers package caches the build; subsequent calls are no-ops.
func buildOnce(t *testing.T) {
	t.Helper()
	force := false
	runPagefind := false
	helpers.Build(t, force, runPagefind)
}

# MiNNaK Hugo Theme — developer tasks
#
# Usage:
#   make dev            — build + index search + serve full demo at http://localhost:1313/
#   make dev-nosearch   — build (no pagefind) + serve at http://localhost:1313/
#   make watch          — live-reload dev server via hugo server (no search)
#   make build          — build the example site (hugo only)
#   make pagefind       — build the example site + run pagefind indexer
#   make test           — run Go markup tests (reuses existing public/ if present)
#   make e2e            — run Playwright E2E tests (starts hugo server automatically)
#   make e2e-search     — run Playwright + search tests (requires make pagefind first)
#   make lint           — hugo path-warnings lint pass
#   make ci             — full CI sequence: build + pagefind + lint + test + e2e-search
#   make clean          — remove built output

.PHONY: build pagefind dev dev-nosearch watch test test-fresh e2e e2e-search e2e-ui lint ci clean

EXAMPLE_SITE := exampleSite
PUBLIC_DIR   := $(EXAMPLE_SITE)/public
DEV_PORT     := 1313

# Full demo: build, index search, then serve the static output.
# Consumers control search with params.search = false in their hugo.toml.
dev: pagefind
	@printf '\n  Demo site → http://localhost:$(DEV_PORT)/\n  Ctrl-C to stop\n\n'
	python3 -m http.server $(DEV_PORT) --bind 127.0.0.1 --directory $(PUBLIC_DIR)

# Same as dev but without the pagefind index or search UI.
# Useful when pagefind is unavailable or you want a faster iteration cycle.
dev-nosearch:
	HUGO_PARAMS_SEARCH=false hugo --source $(EXAMPLE_SITE) --logLevel warn
	@printf '\n  Demo site (no search) → http://localhost:$(DEV_PORT)/\n  Ctrl-C to stop\n\n'
	python3 -m http.server $(DEV_PORT) --bind 127.0.0.1 --directory $(PUBLIC_DIR)

# Fast live-reload development server.
# Search is disabled so the browser makes no requests for missing pagefind files.
watch:
	HUGO_PARAMS_SEARCH=false hugo server --source $(EXAMPLE_SITE) --port $(DEV_PORT) \
	    --disableFastRender --navigateToChanged

build:
	hugo --source $(EXAMPLE_SITE) --logLevel warn

pagefind: build
	npx --yes pagefind@latest --site $(PUBLIC_DIR)

lint:
	hugo --source $(EXAMPLE_SITE) --renderToMemory \
	     --printPathWarnings --logLevel warn

test:
	cd tests && MINNAK_SKIP_BUILD=1 go test ./... -count=1 -v

# Run tests and rebuild first (used when you haven't built yet)
test-fresh:
	cd tests && go test ./... -count=1 -v

e2e:
	cd tests/e2e && npx playwright test

e2e-search:
	cd tests/e2e && MINNAK_RUN_PAGEFIND=1 npx playwright test

e2e-ui:
	cd tests/e2e && npx playwright test --ui

ci: build pagefind lint test e2e-search

clean:
	rm -rf $(PUBLIC_DIR)
	rm -rf tests/e2e/playwright-report
	rm -rf tests/e2e/test-results

# MiNNaK Hugo Theme — developer tasks
#
# Usage:
#   make build        — build the example site (hugo only)
#   make pagefind     — build the example site + run pagefind indexer
#   make test         — run Go markup tests (reuses existing public/ if present)
#   make e2e          — run Playwright E2E tests (starts hugo server automatically)
#   make ci           — full CI sequence: build + lint + test + e2e
#   make clean        — remove built output

.PHONY: build pagefind test e2e lint ci clean

EXAMPLE_SITE := exampleSite
PUBLIC_DIR   := $(EXAMPLE_SITE)/public

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

e2e-ui:
	cd tests/e2e && npx playwright test --ui

ci: build lint test e2e

clean:
	rm -rf $(PUBLIC_DIR)
	rm -rf tests/e2e/playwright-report
	rm -rf tests/e2e/test-results

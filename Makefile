# MiNNaK Hugo Theme — developer tasks

.PHONY: help build pagefind dev dev-nosearch watch test test-fresh e2e e2e-search e2e-ui lint ci clean

EXAMPLE_SITE := exampleSite
PUBLIC_DIR   := $(EXAMPLE_SITE)/public
DEV_PORT     := 1313

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
	     /^[a-zA-Z0-9_-]+:.*##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@printf '\n'

dev: pagefind ## Build + index search + serve demo at http://localhost:$(DEV_PORT)/
	@printf '\n  Demo site → http://localhost:$(DEV_PORT)/\n  Ctrl-C to stop\n\n'
	python3 -m http.server $(DEV_PORT) --bind 127.0.0.1 --directory $(PUBLIC_DIR)

dev-nosearch: ## Build (no pagefind) + serve at http://localhost:$(DEV_PORT)/
	HUGO_PARAMS_SEARCH=false hugo --source $(EXAMPLE_SITE) --logLevel warn
	@printf '\n  Demo site (no search) → http://localhost:$(DEV_PORT)/\n  Ctrl-C to stop\n\n'
	python3 -m http.server $(DEV_PORT) --bind 127.0.0.1 --directory $(PUBLIC_DIR)

watch: ## Live-reload dev server via hugo server (search disabled)
	HUGO_PARAMS_SEARCH=false hugo server --source $(EXAMPLE_SITE) --port $(DEV_PORT) \
	    --disableFastRender --navigateToChanged

build: ## Build the example site (hugo only)
	hugo --source $(EXAMPLE_SITE) --logLevel warn

pagefind: build ## Build example site + run pagefind indexer
	npx --yes pagefind@latest --site $(PUBLIC_DIR)

lint: ## Hugo path-warnings lint pass
	hugo --source $(EXAMPLE_SITE) --renderToMemory \
	     --printPathWarnings --logLevel warn

test: ## Run Go markup tests (reuses existing public/ if present)
	cd tests && MINNAK_SKIP_BUILD=1 go test ./... -count=1 -v

test-fresh: ## Run Go tests after a fresh build
	cd tests && go test ./... -count=1 -v

e2e: ## Run Playwright E2E tests
	cd tests/e2e && npx playwright test

e2e-search: ## Run Playwright + search tests (requires make pagefind first)
	cd tests/e2e && MINNAK_RUN_PAGEFIND=1 npx playwright test

e2e-ui: ## Open Playwright UI mode
	cd tests/e2e && npx playwright test --ui

ci: build pagefind lint test e2e-search ## Full CI sequence

clean: ## Remove built output and test artifacts
	rm -rf $(PUBLIC_DIR)
	rm -rf tests/e2e/playwright-report
	rm -rf tests/e2e/test-results

# MiNNaK Hugo Theme

```
   /\_/\
  ( ^.^ )
   __ __  _  _ _  _ _  ___  _ __
  |  \  \| || \ || \ || . || / /
  |     || ||   ||   ||   ||  \
  |_|_|_||_||_\_||_\_||_|_||_\_\
```

A port of the [MiNNaK Ghost theme](https://github.com/intentionally-left-nil/minnak-hugo) to [Hugo](https://gohugo.io/), styled after [terminal.space](https://terminal.space). Minimal, dark, responsive.

## Requirements

- Hugo **extended** ≥ 0.120.0
- Go ≥ 1.21 (for Hugo module resolution)

## Install as a Hugo module

**1. Initialise your site as a module** (skip if you already have a `go.mod`):

```bash
hugo mod init github.com/youruser/yoursite
```

**2. Add the theme to your `hugo.toml`:**

```toml
[module]
  [[module.imports]]
    path = "github.com/intentionally-left-nil/minnak-hugo"
```

**3. Pull the module:**

```bash
hugo mod get github.com/intentionally-left-nil/minnak-hugo
```

## Configuration

Minimal `hugo.toml` for a new site:

```toml
baseURL  = "https://yoursite.example/"
title    = "Your Site Title"

[taxonomies]
  category = "categories"
  tag      = "tags"

paginate = 15

[params]
  description = "Your site description"
  search      = true          # set false to disable Pagefind mount

# Main navigation — drives the Categories tab in the sidebar
[menu]
  [[menu.main]]
    name   = "Technology"
    url    = "/categories/technology/"
    weight = 1
  [[menu.main]]
    name   = "AI"
    url    = "/categories/ai/"
    weight = 2
```

## Content

Posts go in `content/posts/`. Front matter convention:

```yaml
---
title: "Post Title"
date: 2026-03-15T10:00:00Z
categories: ["Technology"]   # exactly one — drives the sidebar category tab
tags: ["rust", "systems"]    # optional; shown as pills in post footer
summary: "Short excerpt shown on cards."
---
```

Standalone pages (About, etc.) go anywhere else in `content/` with `type: page`:

```yaml
---
title: "About"
type: "page"
---
```

Category term pages can have a description via `content/categories/<slug>/_index.md`:

```yaml
---
title: "Technology"
description: "Computer technology and coding articles"
---
```

## Search (Pagefind)

The theme mounts [Pagefind](https://pagefind.app/) in the sidebar's search tab. Pagefind runs as a post-build step — it is **not** needed during `hugo server` development.

**Build with search index:**

```bash
hugo                                    # build the site
npx --yes pagefind@latest --site public # index the output
```

Or use the example site Makefile:

```bash
make pagefind
```

During `hugo server` development the `#search` div is present but the Pagefind UI gracefully shows a disabled state until an index is available.

To disable search entirely, set `params.search = false` in your config.

## Fonts

The theme ships subset LatoLatin font files (`static/fonts/`). They are served at `/fonts/` and referenced by the bundled CSS — no extra steps needed.

## Development (this repo)

```bash
# Build the example site
make build

# Build + index search
make pagefind

# Run Go markup tests (50 assertions)
make test

# Run Playwright E2E tests (sidebar JS, responsive layout)
make e2e

# Run Playwright tests including search (requires make pagefind first)
make e2e-search

# Full CI sequence (what GitHub Actions runs)
make ci

# Remove build output
make clean
```

### TDD workflow

Each feature has a corresponding test in `tests/`. The cycle is:

1. Write a failing assertion in `tests/*_test.go`
2. Run `make test` — see it fail
3. Edit layouts / partials / CSS in the module root
4. Re-run `make test` — see it pass
5. Commit

For JS/CSS behaviour, use `make e2e` with the Playwright specs in `tests/e2e/specs/`.

## License

GNU General Public License v2 or later — <http://www.gnu.org/licenses/gpl-2.0.html>

This theme is a port of MiNNaK by Tamer Mancar, distributed under GNU GPL v2 or later.

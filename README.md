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

# The theme expects singular taxonomy keys on both sides — front matter
# uses `category:` / `tag:` (singular) and the rendered URLs are
# /category/<slug>/ and /tag/<slug>/. This matches the convention used
# by sites migrated from WordPress.
[taxonomies]
  category = "category"
  tag      = "tag"

paginate = 15

[params]
  description = "Your site description"
  search      = true          # set false to disable Pagefind mount

# RSS: generate a single site-wide feed at /rss.xml.
# Without this, Hugo also emits feeds for every tag and category page.
[outputs]
  home     = ['html', 'rss']
  section  = ['html']
  taxonomy = ['html']
  term     = ['html']

[outputFormats]
  [outputFormats.RSS]
    baseName = "rss"          # publish to /rss.xml instead of /index.xml

# Main navigation — drives the Categories tab in the sidebar.
# URLs use the singular /category/ form to match the [taxonomies] block.
[menu]
  [[menu.main]]
    name   = "Technology"
    url    = "/category/technology/"
    weight = 1
  [[menu.main]]
    name   = "AI"
    url    = "/category/ai/"
    weight = 2
```

## Content

Posts go in `content/posts/`. Front matter convention:

```yaml
---
title: "Post Title"
date: 2026-03-15T10:00:00Z
category: ["Technology"]    # exactly one — drives the sidebar category tab
tag: ["rust", "systems"]    # optional; shown as pills in post footer
summary: "Short excerpt shown on cards."
---
```

> **Note on singular keys:** the theme reads `.Params.category` and
> `.Params.tag`, and links to `/category/<slug>/` and `/tag/<slug>/`.
> Make sure the front matter keys (and the `[taxonomies]` block in
> `hugo.toml`) are singular — using plural `categories:` / `tags:` will
> result in empty taxonomy widgets and broken links.

### Guest author (optional)

Posts can declare a guest author via an optional `author:` key. When
set, the post header renders a `Guest author: NAME` byline and `<head>`
emits SEO metadata: `<meta name="author">`, `<meta property="article:author">`,
and a [schema.org `Article`](https://schema.org/Article) JSON-LD block.
Posts without `author:` render unchanged — older posts and main-author
posts are unaffected.

Two shapes are accepted. Plain string:

```yaml
---
title: "Guest Post"
author: "Jane Doe"
---
```

Map (use when the author has a homepage):

```yaml
---
title: "Guest Post"
author:
  name: "Jane Doe"
  url:  "https://janedoe.com"   # optional; wraps the byline name in <a>
  email: "jane@example.com"     # optional; used by RSS <author>
---
```

Without `email`, the RSS `<author>` element is omitted for that item
(RSS spec requires an email; falling back to the site owner would
misattribute the post).

Standalone pages (About, etc.) go anywhere else in `content/` with `type: page`:

```yaml
---
title: "About"
type: "page"
---
```

Category term pages can have a description via `content/category/<slug>/_index.md`:

```yaml
---
title: "Technology"
description: "Computer technology and coding articles"
---
```

### Cover / feature images

Each post can have a hero image rendered on its summary card. The theme
resolves it in this order:

1. **`feature.*`** in the page bundle (preferred for new content):

   ```text
   content/posts/my-post/
   ├── index.md
   └── feature.jpg
   ```

   No front-matter changes needed — `feature.jpg`, `feature.png`,
   `feature.webp`, etc. are all picked up automatically.

2. **`cover.image`** in front matter (for content migrated from
   WordPress, where images live under `images/`):

   ```yaml
   ---
   title: "Two Years With a Fairphone"
   cover:
     image: images/Fairphone.jpg
     alt:   "A Fairphone on a wooden desk"
   ---
   ```

   The path is resolved relative to the page bundle, so the image must
   exist at `content/posts/my-post/images/Fairphone.jpg`.

In either case the theme generates a `<picture>` block with WebP and
JPEG sources at six widths (30, 100, 300, 600, 1200, 2000 px). The alt
text falls back from `cover.alt` → page title.

### Gallery shortcode

Drop multiple images into a page bundle under `images/gallery/` and
declare captions in the post's `resources:` front matter:

```yaml
---
title: "Mount Rainier in Four Frames"
date: 2026-04-01T08:00:00Z
category: ["Photography"]
resources:
  - src: "images/gallery/photo-01.jpg"
    title: "Sunrise on the eastern flank"
  - src: "images/gallery/photo-02.jpg"
    title: "Tipsoo Lake reflection at dawn"
---

Some prose…

{{</* gallery cols="3" */>}}
```

The shortcode emits a `<figure>` per image with a square thumbnail
(cropped via Hugo's `.Fill`) linked to the full-size original. The
caption comes from each resource's `.Title`. Images are enumerated in
the order they appear in the `resources:` block.

The `cols` parameter (default `1`) controls the desktop column count.
On smaller viewports the grid collapses automatically:

| Viewport          | `cols="3"` actual columns |
|-------------------|---------------------------|
| ≥ 992 px (desktop)| 3                         |
| 601–991 px (tablet)| 2                        |
| ≤ 600 px (mobile) | 1                         |

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
# Full demo: build + index search + serve at http://localhost:1313/
make dev

# Fast live-reload dev server via hugo server (search gracefully absent)
make watch

# Build the example site (hugo only)
make build

# Build + index search (outputs to exampleSite/public/pagefind/)
make pagefind

# Run Go markup tests
make test

# Run Playwright E2E tests (sidebar JS, responsive layout, gallery)
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

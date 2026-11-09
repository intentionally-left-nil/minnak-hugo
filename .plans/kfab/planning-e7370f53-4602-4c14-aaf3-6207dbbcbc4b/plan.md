<!-- kfab-plan schema=1 -->
# Implementation plan

## Change title

Remove obsolete pagefind same-page hashchange search modal workaround

## Problem

assets/js/search.js closes the pagefind modal on hashchange, a workaround for pagefind issue 1125. The build uses pagefind@latest, and the request says to assume the upstream fix is included, so the workaround is now unnecessary.

## Evidence

- assets/js/search.js default-exports fixSamePageSearch, which calls modal.close() on window hashchange and is marked as a workaround for pagefind issue 1125.
- assets/js/index.js imports and calls fixSamePageSearch. Triage found no other reference to it.
- The Makefile and CI use pagefind@latest, so there is no pinned version to bump.
- Triage found no hashchange references in tests.

## Desired outcome

The workaround code is gone and the theme relies on pagefind's own behaviour for closing the modal.

## Scope

- Delete assets/js/search.js.
- Remove the fixSamePageSearch import and call from assets/js/index.js.

## Exclusions

- No changes to pagefind versioning, the Makefile, CI or the README.
- No changes to sidebar or lightbox code.

## Approach

- Delete assets/js/search.js.
- Edit assets/js/index.js to drop the import and the call, leaving sidebar() and lightbox() unchanged.
- Search the repository for any remaining reference to search.js or fixSamePageSearch.

## Affected areas

- assets/js/search.js
- assets/js/index.js

## Acceptance criteria

### ac-1: assets/js/search.js no longer exists, and nothing in the repository references fixSamePageSearch or imports search.js. index.js still initialises sidebar and lightbox.

**Verification:** Source review

### ac-2: With the latest pagefind, clicking a same-page anchor sub-result in the search modal closes the modal.

**Verification:** Manual: build the site with pagefind@latest, search, click a same-page heading result, and confirm the modal closes.

## Risks

- The upstream fix was not verified. If it is absent from the latest pagefind, the modal will stay open after same-page navigation.

## Assumptions

- As the request states, pagefind@latest includes the upstream fix for issue 1125.

## Decisions

- None recorded

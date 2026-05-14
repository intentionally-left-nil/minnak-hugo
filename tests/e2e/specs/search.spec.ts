/**
 * search.spec.ts — Pagefind search integration (M2.4)
 *
 * These tests require a pre-built site with a pagefind index.
 * They only run when MINNAK_RUN_PAGEFIND=1 (set in CI and via `make e2e-search`).
 *
 * Tests:
 *  - Search mount exists in the widget tab
 *  - Typing a query produces results
 *  - A result link navigates to the correct post
 *  - The search input is styled with the dark theme (Pagefind CSS overrides)
 */

import { test, expect } from '@playwright/test';
import { runPagefind } from '../playwright.config';

test.describe('Pagefind search', () => {
  test.skip(!runPagefind, 'Pagefind index not built — set MINNAK_RUN_PAGEFIND=1');

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    // Switch to the widget tab to reveal the search input
    await page.locator('#v-pills-widgets-tab').dispatchEvent('mousedown');
    await page.waitForSelector('#search .pagefind-ui__search-input', { timeout: 5000 });
  });

  test('search input is visible in widget tab', async ({ page }) => {
    await expect(page.locator('#search .pagefind-ui__search-input')).toBeVisible();
  });

  test('typing a query returns results', async ({ page }) => {
    await page.locator('#search .pagefind-ui__search-input').fill('rust');
    // Pagefind debounces input; wait for results
    await page.waitForSelector('#search .pagefind-ui__result', { timeout: 5000 });
    const results = page.locator('#search .pagefind-ui__result');
    await expect(results).toHaveCount(1);  // Only one post mentions "rust" prominently
  });

  test('a search result links to the correct post', async ({ page }) => {
    await page.locator('#search .pagefind-ui__search-input').fill('diffusion');
    await page.waitForSelector('#search .pagefind-ui__result-title a', { timeout: 5000 });
    const firstResultLink = page.locator('#search .pagefind-ui__result-title a').first();
    await expect(firstResultLink).toHaveAttribute('href', /diffusion-models/);
  });

  test('clicking a search result navigates to the post', async ({ page }) => {
    await page.locator('#search .pagefind-ui__search-input').fill('cantonese');
    await page.waitForSelector('#search .pagefind-ui__result-title a', { timeout: 5000 });
    const firstLink = page.locator('#search .pagefind-ui__result-title a').first();
    const href = await firstLink.getAttribute('href');
    await firstLink.click();
    await expect(page).toHaveURL(new RegExp('cantonese'));
  });

  test('dark theme CSS overrides are applied to pagefind UI', async ({ page }) => {
    const input = page.locator('#search .pagefind-ui__search-input');
    const bgColor = await input.evaluate(el => getComputedStyle(el).backgroundColor);
    // Should be var(--bg-page) = #1A1C1E = rgb(26, 28, 30)
    expect(bgColor).toBe('rgb(26, 28, 30)');
  });
});

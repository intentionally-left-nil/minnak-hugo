/**
 * search.spec.ts — Pagefind Component UI search integration
 *
 * These tests require a pre-built site with a pagefind index.
 * They only run when MINNAK_RUN_PAGEFIND=1 (set in CI and via `make e2e-search`).
 *
 * Tests:
 *  - Modal trigger is visible in the sidebar header without any tab interaction
 *  - Clicking the trigger opens the search modal
 *  - Typing a query inside the modal produces results
 *  - A result link points to the correct post
 *  - Pressing Escape closes the modal
 *  - Ctrl+K keyboard shortcut opens the modal
 *
 * Note: <pagefind-modal> renders its <dialog> inside shadow DOM via the top
 * layer, so Playwright's toBeVisible() always sees the host element as hidden.
 * We use toHaveJSProperty('isOpen', ...) against the component's public API
 * to assert open/closed state.
 */

import { test, expect } from '@playwright/test';
import { runPagefind } from '../playwright.config';

test.describe('Pagefind Component UI search', () => {
  test.skip(!runPagefind, 'Pagefind index not built — set MINNAK_RUN_PAGEFIND=1');

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('search trigger is visible in the sidebar header', async ({ page }) => {
    await expect(page.locator('pagefind-modal-trigger')).toBeVisible();
  });

  test('clicking the trigger opens the search modal', async ({ page }) => {
    await page.locator('pagefind-modal-trigger').click();
    await expect(page.locator('pagefind-modal')).toHaveJSProperty('isOpen', true);
  });

  test('Ctrl+K opens the search modal', async ({ page }) => {
    await page.keyboard.press('Control+k');
    await expect(page.locator('pagefind-modal')).toHaveJSProperty('isOpen', true);
  });

  test('Escape closes the search modal', async ({ page }) => {
    await page.locator('pagefind-modal-trigger').click();
    await expect(page.locator('pagefind-modal')).toHaveJSProperty('isOpen', true);
    await page.keyboard.press('Escape');
    await expect(page.locator('pagefind-modal')).toHaveJSProperty('isOpen', false);
  });

  test('typing a query returns results', async ({ page }) => {
    await page.locator('pagefind-modal-trigger').click();
    const input = page.locator('pagefind-modal pagefind-input').getByRole('searchbox');
    await input.fill('rust');
    await page.waitForSelector('pagefind-modal pagefind-results li', { timeout: 5000 });
    const results = page.locator('pagefind-modal pagefind-results li');
    await expect(results).not.toHaveCount(0);
  });

  test('a search result links to the correct post', async ({ page }) => {
    await page.locator('pagefind-modal-trigger').click();
    const input = page.locator('pagefind-modal pagefind-input').getByRole('searchbox');
    await input.fill('diffusion');
    await page.waitForSelector('pagefind-modal pagefind-results a', { timeout: 5000 });
    const firstLink = page.locator('pagefind-modal pagefind-results a').first();
    await expect(firstLink).toHaveAttribute('href', /diffusion-models/);
  });
});

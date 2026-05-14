/**
 * responsive.spec.ts — M2.6 Sidebar CSS + M5 Grid responsive visibility
 *
 * Tests:
 *  - Desktop (≥992px): sidebar visible, 3-col grid shown, 2-col/1-col hidden
 *  - Tablet (600–991px): sidebar hidden off-screen, 2-col shown, 3-col hidden
 *  - Mobile (<600px): 1-col shown, 2-col and 3-col hidden
 *  - Desktop: toggle button hidden
 *  - Mobile: toggle button visible
 */

import { test, expect } from '@playwright/test';

// ─── Desktop layout (1280×800) ──────────────────────────────────────────────

test.describe('Desktop layout (1280×800)', () => {
  test.use({ viewport: { width: 1280, height: 800 } });

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('sidebar is visible on desktop', async ({ page }) => {
    const sidebar = page.locator('#left-sidebar');
    const box = await sidebar.boundingBox();
    expect(box).not.toBeNull();
    // Sidebar left edge should be at x≥0 (on screen)
    expect(box!.x).toBeGreaterThanOrEqual(0);
  });

  test('3-col grid is visible on desktop', async ({ page }) => {
    await expect(page.locator('.gh-postfeed-3col')).toBeVisible();
  });

  test('2-col grid is hidden on desktop', async ({ page }) => {
    await expect(page.locator('.gh-postfeed-2col')).toBeHidden();
  });

  test('1-col grid is hidden on desktop', async ({ page }) => {
    await expect(page.locator('.gh-postfeed-1col')).toBeHidden();
  });

  test('toggle button is not visible on desktop', async ({ page }) => {
    await expect(page.locator('#toggle-button')).toBeHidden();
  });
});

// ─── Tablet layout (768×1024) ───────────────────────────────────────────────

test.describe('Tablet layout (768×1024)', () => {
  test.use({ viewport: { width: 768, height: 1024 } });

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('sidebar is off-screen on tablet (left < 0)', async ({ page }) => {
    const sidebar = page.locator('#left-sidebar');
    const box = await sidebar.boundingBox();
    expect(box).not.toBeNull();
    // Sidebar should be translated off the left edge
    expect(box!.x + box!.width).toBeLessThanOrEqual(0);
  });

  test('3-col grid is hidden on tablet', async ({ page }) => {
    await expect(page.locator('.gh-postfeed-3col')).toBeHidden();
  });

  test('2-col grid is visible on tablet', async ({ page }) => {
    await expect(page.locator('.gh-postfeed-2col')).toBeVisible();
  });

  test('1-col grid is hidden on tablet', async ({ page }) => {
    await expect(page.locator('.gh-postfeed-1col')).toBeHidden();
  });

  test('toggle button is visible on tablet', async ({ page }) => {
    await expect(page.locator('#toggle-button')).toBeVisible();
  });
});

// ─── Mobile layout (390×844) ────────────────────────────────────────────────

test.describe('Mobile layout (390×844)', () => {
  test.use({ viewport: { width: 390, height: 844 } });

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('3-col grid is hidden on mobile', async ({ page }) => {
    await expect(page.locator('.gh-postfeed-3col')).toBeHidden();
  });

  test('2-col grid is hidden on mobile', async ({ page }) => {
    await expect(page.locator('.gh-postfeed-2col')).toBeHidden();
  });

  test('1-col grid is visible on mobile', async ({ page }) => {
    await expect(page.locator('.gh-postfeed-1col')).toBeVisible();
  });
});

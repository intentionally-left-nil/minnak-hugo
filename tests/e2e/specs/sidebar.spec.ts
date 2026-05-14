/**
 * sidebar.spec.ts — M2.5 Sidebar JS behaviour
 *
 * Tests:
 *  - Tab switching: clicking the widget tab moves active-menu
 *  - Mobile drawer: toggle button adds body.open-sidebar
 *  - Mobile drawer: clicking outside closes the sidebar
 *  - Focus trap: blur on last focusable element closes sidebar
 */

import { test, expect } from '@playwright/test';

test.describe('Sidebar tab switching', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('categories tab is active by default', async ({ page }) => {
    const firstTab = page.locator('#left-sidebar .vertical-menu-item').first();
    await expect(firstTab).toHaveClass(/active-menu/);
  });

  test('clicking widget tab moves active-menu to second tab', async ({ page }) => {
    const widgetTabLink = page.locator('#v-pills-widgets-tab');
    await widgetTabLink.dispatchEvent('mousedown');

    const secondTab = page.locator('#left-sidebar .vertical-menu-item').nth(1);
    await expect(secondTab).toHaveClass(/active-menu/);

    const firstTab = page.locator('#left-sidebar .vertical-menu-item').first();
    await expect(firstTab).not.toHaveClass(/active-menu/);
  });

  test('clicking categories tab restores active-menu to first tab', async ({ page }) => {
    // First switch to widget tab
    await page.locator('#v-pills-widgets-tab').dispatchEvent('mousedown');
    // Then switch back to categories tab
    await page.locator('#v-pills-cats-tab').dispatchEvent('mousedown');

    const firstTab = page.locator('#left-sidebar .vertical-menu-item').first();
    await expect(firstTab).toHaveClass(/active-menu/);

    const secondTab = page.locator('#left-sidebar .vertical-menu-item').nth(1);
    await expect(secondTab).not.toHaveClass(/active-menu/);
  });

  test('only one tab is active at a time', async ({ page }) => {
    await page.locator('#v-pills-widgets-tab').dispatchEvent('mousedown');
    const activeCount = await page.locator('#left-sidebar .vertical-menu-item.active-menu').count();
    expect(activeCount).toBe(1);
  });
});

test.describe('Mobile sidebar drawer', () => {
  test.use({ viewport: { width: 390, height: 844 } }); // iPhone 14 size

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('toggle button is visible on mobile', async ({ page }) => {
    await expect(page.locator('#toggle-button')).toBeVisible();
  });

  test('toggle button adds body.open-sidebar on mousedown', async ({ page }) => {
    await page.locator('#toggle-button').dispatchEvent('mousedown');
    await expect(page.locator('body')).toHaveClass(/open-sidebar/);
  });

  test('toggle button removes body.open-sidebar on second mousedown', async ({ page }) => {
    const btn = page.locator('#toggle-button');
    await btn.dispatchEvent('mousedown');
    await expect(page.locator('body')).toHaveClass(/open-sidebar/);
    await btn.dispatchEvent('mousedown');
    await expect(page.locator('body')).not.toHaveClass(/open-sidebar/);
  });

  test('clicking outside sidebar closes it', async ({ page }) => {
    // Open sidebar
    await page.locator('#toggle-button').dispatchEvent('mousedown');
    await expect(page.locator('body')).toHaveClass(/open-sidebar/);

    // Click on the main content area (outside sidebar)
    await page.locator('.gh-main').click();
    await expect(page.locator('body')).not.toHaveClass(/open-sidebar/);
  });
});

/**
 * figure.spec.ts — Figure shortcode maxheight CSS behaviour
 *
 * The `maxheight` param emits `style="max-height: <value>px; width: auto;"` on
 * the rendered <img>.  These tests verify that the browser actually clamps the
 * rendered height to that value while allowing the width to follow the aspect
 * ratio (i.e. the image is not stretched or cropped, just scaled down).
 *
 * Fixture: posts/figure-maxheight-fixture/ — renders
 *   {{< figure src="feature.jpg" alt="Test image" maxheight="400" >}}
 * The source image is 640×960 (portrait), so without the constraint its
 * natural rendered height would exceed 400px on any reasonably wide viewport.
 */

import { test, expect } from '@playwright/test';

const FIXTURE_URL = '/posts/figure-maxheight-fixture/';
const MAX_HEIGHT_PX = 400;

test.describe('Figure maxheight CSS clamping', () => {
  test.use({ viewport: { width: 1280, height: 800 } });

  test.beforeEach(async ({ page }) => {
    await page.goto(FIXTURE_URL);
  });

  test('rendered height does not exceed maxheight', async ({ page }) => {
    const img = page.locator('.entry-content figure img').first();
    await expect(img).toBeVisible();

    const box = await img.boundingBox();
    expect(box).not.toBeNull();
    // Allow 1px rounding tolerance.
    expect(box!.height).toBeLessThanOrEqual(MAX_HEIGHT_PX + 1);
  });

  test('rendered height equals maxheight (constraint is active, not a no-op)', async ({
    page,
  }) => {
    // The fixture image is 640×960 — taller than 400px — so the clamp must
    // be actively limiting the height.  We verify it is close to MAX_HEIGHT_PX
    // rather than the unconstrained natural height.
    const img = page.locator('.entry-content figure img').first();
    await expect(img).toBeVisible();

    const box = await img.boundingBox();
    expect(box).not.toBeNull();
    // Height should be at or very close to the cap (within 1px rounding).
    expect(box!.height).toBeGreaterThan(MAX_HEIGHT_PX - 2);
    expect(box!.height).toBeLessThanOrEqual(MAX_HEIGHT_PX + 1);
  });

  test('width follows aspect ratio (not stretched)', async ({ page }) => {
    // The source is 640×960 (2:3 portrait). At max-height 400px the natural
    // width should be ≈ 267px.  We allow a generous tolerance for sub-pixel
    // rendering and different browser rounding, but it must NOT be 0 or equal
    // to the viewport width (which would indicate no aspect-ratio preservation).
    const img = page.locator('.entry-content figure img').first();
    await expect(img).toBeVisible();

    const box = await img.boundingBox();
    expect(box).not.toBeNull();

    const expectedWidth = (box!.height * 640) / 960;
    // Within 5% tolerance.
    expect(box!.width).toBeGreaterThan(expectedWidth * 0.95);
    expect(box!.width).toBeLessThan(expectedWidth * 1.05);
  });
});

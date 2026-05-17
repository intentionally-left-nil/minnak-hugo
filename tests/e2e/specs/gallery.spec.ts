/**
 * gallery.spec.ts — Gallery shortcode responsive layout
 *
 * The shortcode emits <div class="gallery gallery-cols-N">…</div>. Theme
 * CSS caps the column count on smaller viewports:
 *
 *   Desktop (≥992px):           gallery-cols-N renders as N columns
 *   Tablet  (601–991px):        gallery-cols-{3,4} fall back to 2 columns;
 *                               gallery-cols-2 keeps 2 columns
 *   Mobile  (≤600px):           every gallery-cols-N collapses to 1 column
 *
 * The fixture post posts/gallery-mount-rainier/ uses cols="3", so we test
 * the 3-col case across all three breakpoints.
 */

import { test, expect, Locator } from '@playwright/test';

const GALLERY_URL = '/posts/gallery-mount-rainier/';

/**
 * Compute the rendered column count by reading grid-template-columns.
 * The browser resolves `repeat(N, minmax(0, 1fr))` to N space-separated
 * pixel values, so the count of pixel tokens is the column count.
 */
async function renderedColumnCount(gallery: Locator): Promise<number> {
  const cols = await gallery.evaluate((el) =>
    window.getComputedStyle(el).gridTemplateColumns,
  );
  return cols.split(/\s+/).filter((s) => s.endsWith('px')).length;
}

/**
 * Group items by their visual row (rounded top y). A grid row may span
 * different heights per item (because of caption length variance), so
 * we round to the nearest 5px to be tolerant.
 */
async function rowCount(items: Locator): Promise<number> {
  const tops = await items.evaluateAll((els) =>
    els.map((el) => Math.round(el.getBoundingClientRect().top / 5) * 5),
  );
  return new Set(tops).size;
}

// ─── Desktop (1280×800) ─────────────────────────────────────────────────────

test.describe('Gallery on desktop (1280×800)', () => {
  test.use({ viewport: { width: 1280, height: 800 } });

  test.beforeEach(async ({ page }) => {
    await page.goto(GALLERY_URL);
  });

  test('cols="3" renders 3 columns', async ({ page }) => {
    const gallery = page.locator('.gallery.gallery-cols-3');
    await expect(gallery).toBeVisible();
    expect(await renderedColumnCount(gallery)).toBe(3);
  });

  test('4 items fit in 2 rows (3 + 1 wrapped)', async ({ page }) => {
    const items = page.locator('.gallery .gallery-item');
    await expect(items).toHaveCount(4);
    expect(await rowCount(items)).toBe(2);
  });

  test('thumbnails are roughly square', async ({ page }) => {
    const firstImg = page.locator('.gallery .gallery-item img').first();
    const box = await firstImg.boundingBox();
    expect(box).not.toBeNull();
    // The actual rendered size depends on column width, but the source
    // is a 400x400 crop so the aspect ratio should be ≈ 1:1.
    const ratio = box!.width / box!.height;
    expect(ratio).toBeGreaterThan(0.95);
    expect(ratio).toBeLessThan(1.05);
  });

  test('clicking a thumbnail links to the full-size original', async ({
    page,
  }) => {
    const firstLink = page.locator('.gallery .gallery-item a').first();
    const href = await firstLink.getAttribute('href');
    expect(href).toMatch(
      /\/posts\/gallery-mount-rainier\/images\/gallery\/photo-01\.jpg$/,
    );
  });
});

// ─── Tablet (768×1024) ──────────────────────────────────────────────────────

test.describe('Gallery on tablet (768×1024)', () => {
  test.use({ viewport: { width: 768, height: 1024 } });

  test.beforeEach(async ({ page }) => {
    await page.goto(GALLERY_URL);
  });

  test('cols="3" collapses to 2 columns at tablet width', async ({ page }) => {
    const gallery = page.locator('.gallery.gallery-cols-3');
    await expect(gallery).toBeVisible();
    expect(await renderedColumnCount(gallery)).toBe(2);
  });

  test('4 items lay out in 2 rows of 2', async ({ page }) => {
    const items = page.locator('.gallery .gallery-item');
    await expect(items).toHaveCount(4);
    expect(await rowCount(items)).toBe(2);
  });

  test('thumbnails do not overflow the viewport horizontally', async ({
    page,
  }) => {
    const items = page.locator('.gallery .gallery-item');
    const overflow = await items.evaluateAll((els) =>
      els.some((el) => {
        const r = el.getBoundingClientRect();
        return r.right > window.innerWidth + 1;
      }),
    );
    expect(overflow).toBe(false);
  });
});

// ─── Mobile (390×844) ───────────────────────────────────────────────────────

test.describe('Gallery on mobile (390×844)', () => {
  test.use({ viewport: { width: 390, height: 844 } });

  test.beforeEach(async ({ page }) => {
    await page.goto(GALLERY_URL);
  });

  test('cols="3" collapses to 1 column at mobile width', async ({ page }) => {
    const gallery = page.locator('.gallery.gallery-cols-3');
    await expect(gallery).toBeVisible();
    expect(await renderedColumnCount(gallery)).toBe(1);
  });

  test('every item occupies its own row (4 rows total)', async ({ page }) => {
    const items = page.locator('.gallery .gallery-item');
    await expect(items).toHaveCount(4);
    expect(await rowCount(items)).toBe(4);
  });

  test('thumbnails fit within the mobile viewport width', async ({ page }) => {
    const firstImg = page.locator('.gallery .gallery-item img').first();
    const box = await firstImg.boundingBox();
    expect(box).not.toBeNull();
    expect(box!.width).toBeLessThanOrEqual(390);
    // And the image should actually be wide enough to be useful (not e.g.
    // squashed to 50px by some flex glitch).
    expect(box!.width).toBeGreaterThan(200);
  });

  test('captions render below thumbnails (not beside them)', async ({
    page,
  }) => {
    const firstItem = page.locator('.gallery .gallery-item').first();
    const img = firstItem.locator('img');
    const caption = firstItem.locator('figcaption');

    const imgBox = await img.boundingBox();
    const capBox = await caption.boundingBox();
    expect(imgBox).not.toBeNull();
    expect(capBox).not.toBeNull();
    // The caption's top edge must sit below the image's bottom edge.
    expect(capBox!.y).toBeGreaterThanOrEqual(imgBox!.y + imgBox!.height - 2);
  });

  test('captions are visible (not clipped to zero height)', async ({
    page,
  }) => {
    const captions = page.locator('.gallery .gallery-item figcaption');
    await expect(captions).toHaveCount(4);
    const heights = await captions.evaluateAll((els) =>
      els.map((el) => el.getBoundingClientRect().height),
    );
    for (const h of heights) {
      expect(h).toBeGreaterThan(0);
    }
  });
});

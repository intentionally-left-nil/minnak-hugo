import { defineConfig, devices } from '@playwright/test';
import path from 'path';

// When MINNAK_RUN_PAGEFIND=1 we serve the pre-built public/ dir (which
// contains the pagefind index) via "hugo server --renderStaticToDisk".
// Otherwise we use the normal in-memory hugo server (no pagefind index,
// search E2E tests are skipped).
const runPagefind = process.env.MINNAK_RUN_PAGEFIND === '1';

// Resolve the exampleSite directory relative to this config file.
const repoRoot = path.resolve(__dirname, '..', '..');
const exampleSiteDir = path.join(repoRoot, 'exampleSite');

export default defineConfig({
  testDir: './specs',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI
    ? [['html', { outputFolder: 'playwright-report', open: 'never' }], ['github']]
    : [['html', { outputFolder: 'playwright-report', open: 'never' }]],

  use: {
    baseURL: 'http://localhost:1314',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  webServer: {
    // hugo server renders from memory; Pagefind search will not be available.
    command: `hugo server --source "${exampleSiteDir}" --port 1314 --disableFastRender --disableLiveReload`,
    port: 1314,
    reuseExistingServer: !process.env.CI,
    stdout: 'ignore',
    stderr: 'pipe',
  },
});

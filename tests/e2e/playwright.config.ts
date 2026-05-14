import { defineConfig, devices } from '@playwright/test';
import path from 'path';

// When MINNAK_RUN_PAGEFIND=1 we serve the pre-built public/ dir (which
// contains the pagefind index) via Python's built-in http.server.
// Otherwise we use the normal in-memory hugo server (no pagefind index,
// search E2E tests are automatically skipped).
export const runPagefind = process.env.MINNAK_RUN_PAGEFIND === '1';

const repoRoot   = path.resolve(__dirname, '..', '..');
const exampleDir = path.join(repoRoot, 'exampleSite');
const publicDir  = path.join(exampleDir, 'public');

const PORT = 1314;

export default defineConfig({
  testDir: './specs',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI
    ? [['html', { outputFolder: 'playwright-report', open: 'never' }], ['github']]
    : [['html', { outputFolder: 'playwright-report', open: 'never' }]],

  use: {
    baseURL: `http://localhost:${PORT}`,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  webServer: runPagefind
    ? {
        // Serve the pre-built public/ dir including the pagefind index.
        // Python's http.server correctly serves Hugo's pretty URL output
        // (directory requests are served as index.html).
        command: `python3 -m http.server ${PORT} --bind 127.0.0.1 --directory "${publicDir}"`,
        port: PORT,
        reuseExistingServer: false,
        stdout: 'ignore',
        stderr: 'pipe',
      }
    : {
        // Fast dev server for non-search tests; pagefind index not available.
        command: `hugo server --source "${exampleDir}" --port ${PORT} --disableFastRender --disableLiveReload`,
        port: PORT,
        reuseExistingServer: !process.env.CI,
        stdout: 'ignore',
        stderr: 'pipe',
      },
});

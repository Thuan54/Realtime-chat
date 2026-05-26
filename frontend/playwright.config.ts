import { defineConfig, devices } from '@playwright/test'
import path from 'path'

const BASE_URL = process.env.TEST_BASE_URL || 'http://localhost:80'

export default defineConfig({
  testDir: './e2e',
  outputDir: './test-results',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html', { outputFolder: './test-results/html', open: 'never' }],
    ['junit', { outputFile: './test-results/e2e-junit.xml' }],
  ],
  use: {
    baseURL: BASE_URL,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    headless: true,
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    // Optional: add firefox, webkit for cross-browser coverage later
  ],
  // Docker Compose lifecycle hooks for containerized staging environment
  globalSetup: require.resolve('./e2e/global-setup'),
  globalTeardown: require.resolve('./e2e/global-teardown'),
  // Expect assertions to use accessible roles/text per NFR-05 (a11y)
  expect: {
    toHaveScreenshot: { maxDiffPixels: 100 },
    timeout: 10000,
  },
})
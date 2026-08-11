import path from 'path'
import { devices, defineConfig } from '@playwright/test'

export default defineConfig({
  globalSetup: './global-setup.ts',
  testDir: './tests',
  timeout: 30 * 1000,
  expect: {
    timeout: 5000,
  },
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['list'],
    [
      'html',
      {
        outputFolder: path.resolve(__dirname, 'automation-files', 'playwright-report'),
        open: 'never',
      },
    ],
  ],
  use: {
    actionTimeout: 0,
    headless: true,
    baseURL: 'https://dotagiftx.com/',
    storageState: 'tests/ui/storageState.json',
    trace: process.env.CI ? 'retain-on-failure' : 'on',
    screenshot: process.env.CI ? 'only-on-failure' : 'on',
    video: 'retain-on-failure',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
  ],
})

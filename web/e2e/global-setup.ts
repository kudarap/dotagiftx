import { chromium, FullConfig } from '@playwright/test'

async function globalSetup(config: FullConfig) {
  const storage = 'e2e/storageState.json'
  const browser = await chromium.launch()
  const context = await browser.newContext()
  const page = await context.newPage()

  await page.goto('http://localhost:3000')
  // TODO: add login steps here if the tests require authenticated state.

  await context.storageState({ path: storage })
  await browser.close()
}

export default globalSetup

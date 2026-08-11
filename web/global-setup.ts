import { chromium, FullConfig } from '@playwright/test'

async function globalSetup(config: FullConfig) {
  const storage = 'tests/ui/storageState.json'
  const browser = await chromium.launch()
  const context = await browser.newContext()
  const page = await context.newPage()

  await page.goto('https://dotagiftx.com/')
  // TODO: add login steps here if the tests require authenticated state.

  await context.storageState({ path: storage })
  await browser.close()
}

export default globalSetup

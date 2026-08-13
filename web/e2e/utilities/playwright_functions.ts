import { Locator, Page } from '@playwright/test'
import TimeoutError from '@playwright/test'

export class PlaywrightFunctions {
  public readonly page: Page

  constructor(page: Page) {
    this.page = page
  }

  /**
   * Waits for the page to load completely within the specified timeout.
   * @param page - The Playwright page instance to wait for.
   * @param timeout  - The maximum time to wait for the page to load (in milliseconds). Default is 3000ms.
   * @returns - A promise that resolves to true if the page loads successfully within the timeout, or false if it times out.
   */
  async waitForPageToLoad(page: Page, timeout: number = 10000): Promise<boolean> {
    try {
      await page.waitForLoadState('domcontentloaded', { timeout })
      await page.waitForLoadState('load', { timeout })
      return true
    } catch (error) {
      return false
    }
  }

  /**
   * Navigates to the specified URL using the Playwright page instance.
   * @param url - The URL to navigate to.
   */
  async navigateToUrl(url: string): Promise<void> {
    await this.page.goto(url)
    await this.waitForPageToLoad(this.page)
  }

  async clickThenWait(locator: Locator): Promise<void> {
    await locator.click()
    await this.waitForPageToLoad(this.page)
  }

  /**
   * Waits for the specified locator to become visible within the given timeout.
   * Retries visibility check at the specified polling interval if initial wait fails.
   * @param locator - The Playwright Locator to check for visibility.
   * @param timeout - Maximum time to wait for the locator to become visible (in milliseconds). Default is 5000.
   * @param polling - Interval between visibility checks when retrying (in milliseconds). Default is 100.
   */
  async isLocatorVisible(
    locator: Locator,
    timeout: number = 5000,
    polling: number = 100
  ): Promise<boolean> {
    try {
      await locator.waitFor({ state: 'visible', timeout })
      await locator.scrollIntoViewIfNeeded()
      return true
    } catch (error) {
      console.warn(`Locator ${locator.toString()} not visible, retrying...`)
      const startTime = Date.now()
      while (Date.now() - startTime < timeout) {
        if (await locator.isVisible()) {
          return true
        }
        await this.page.waitForTimeout(polling)
      }
      console.warn(`Locator ${locator.toString()} is still not visible due: ${error}`)
      return false
    }
  }
}

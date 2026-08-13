import { expect, Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class MobilePage extends CommonPageComponent {
  private readonly lnkMobilePage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkMobilePage = page.locator('a[href="/download"]')
  }

  /**
   * Clicks the Mobile link and verifies that the page is loaded by checking the page title.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyMobilePage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    await this.clickThenWait(this.lnkMobilePage)
    await expect(this.page).toHaveTitle(PageTitles.MOBILE_APP_PAGE)
  }
}

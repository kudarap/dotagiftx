import { expect, Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class DotaGiftxPlusPage extends CommonPageComponent {
  private readonly lnkDotaGiftxPlusPage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkDotaGiftxPlusPage = page.locator('a[href="/plus"]')
  }

  /**
   * Clicks the DotaGiftx Plus link and verifies that the page is loaded by checking the page title.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyDotaGiftxPlusPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    await this.clickThenWait(this.lnkDotaGiftxPlusPage)
    await this.waitForPageToLoad(this.page)
    await expect(this.page).toHaveTitle(PageTitles.DOTAGIFTX_PLUS_PAGE)
  }
}

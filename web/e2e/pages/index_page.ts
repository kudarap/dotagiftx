import { expect, Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class IndexPage extends CommonPageComponent {
  private readonly lnkDotaGiftX: Locator
  constructor(page: Page) {
    super(page)
    this.lnkDotaGiftX = page.locator("a[href='/']")
  }

  /**
   * Clicks the DotaGiftX link and verifies that the main page has loaded by checking the page title.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyDotaGiftXMainPageLoad(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    await this.clickThenWait(this.lnkDotaGiftX.first())
    await expect(this.page).toHaveTitle(PageTitles.MAIN_PAGE)
  }
}

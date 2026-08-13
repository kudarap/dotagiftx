import { expect, Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class BansPage extends CommonPageComponent {
  private readonly lnkBansPage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkBansPage = page.locator('a[href="/bans"]')
  }

  /**
   * Clicks the Bans link and verifies that the page is loaded by checking the page title.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyBansPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    await this.clickThenWait(this.lnkBansPage)
    await expect(this.page).toHaveTitle(PageTitles.BANS_PAGE)
  }
}

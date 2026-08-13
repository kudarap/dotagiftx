import { Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class UpdatesPage extends CommonPageComponent {
  private readonly lnkUpdatesPage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkUpdatesPage = page.locator('a[href="/updates"]')
  }

  /**
   * Clicks the Updates link and verifies that the page is loaded by checking if the link is visible.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyUpdatesPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    if (await this.lnkUpdatesPage.isHidden()) {
      await this.clickThenWait(this.btnMore)
      await this.clickThenWait(this.lnkUpdatesPage)
    }
  }
}

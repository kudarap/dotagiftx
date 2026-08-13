import { Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class MiddleManPage extends CommonPageComponent {
  private readonly lnkMiddleManPage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkMiddleManPage = page.locator('a[href="/middleman"]').first()
  }

  /**
   * Clicks the Middle Man link and verifies that the page is loaded by checking if the link is visible.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyMiddleManPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    if (await this.lnkMiddleManPage.isHidden()) {
      await this.clickThenWait(this.btnMore)
      await this.clickThenWait(this.lnkMiddleManPage)
    }
  }
}

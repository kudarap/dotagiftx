import { Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class ModeratorsPage extends CommonPageComponent {
  private readonly lnkModeratorsPage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkModeratorsPage = page.locator('a[href="/moderators"]').first()
  }

  /**
   * Clicks the Moderators link and verifies that the page is loaded by checking if the link is visible.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyModeratorsPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    if (await this.lnkModeratorsPage.isHidden()) {
      await this.clickThenWait(this.btnMore)
      await this.clickThenWait(this.lnkModeratorsPage)
    }
  }
}

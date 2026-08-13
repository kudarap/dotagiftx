import { expect, Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class TreasuresPage extends CommonPageComponent {
  private readonly lnkTreasuresPage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkTreasuresPage = page.locator('a[href="/treasures"]')
  }

  /**
   * Clicks the Treasures link and verifies that the page is loaded by checking the page title.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyTreasuresPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    await this.clickThenWait(this.lnkTreasuresPage)
    await expect(this.page).toHaveTitle(PageTitles.TREASURES_PAGE)
  }
}

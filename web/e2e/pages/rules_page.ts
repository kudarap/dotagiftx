import { expect, Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class RulesPage extends CommonPageComponent {
  private readonly lnkRulesPage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkRulesPage = page.locator('a[href="/rules"]')
  }

  /**
   * Clicks the Rules link and verifies that the page is loaded by checking the page title.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyRulesPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    await this.clickThenWait(this.lnkRulesPage)
    await expect(this.page).toHaveTitle(PageTitles.RULES_PAGE)
  }
}

import { Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class FaqsPage extends CommonPageComponent {
  private readonly lnkFaqsPage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkFaqsPage = page.locator('a[href="/faqs"]').first()
  }

  /**
   * Clicks the FAQs link and verifies that the page is loaded by checking if the link is visible.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyFaqsPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    if (await this.lnkFaqsPage.isHidden()) {
      await this.clickThenWait(this.btnMore)
      await this.clickThenWait(this.lnkFaqsPage)
    }
  }
}

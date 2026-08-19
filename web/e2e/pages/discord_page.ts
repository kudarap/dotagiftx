import { Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'

export class DiscordPage extends CommonPageComponent {
  private readonly lnkDiscord: Locator
  constructor(page: Page) {
    super(page)
    this.lnkDiscord = page.getByRole('link', { name: 'Discord' })
  }

  /**
   * Clicks the Discord link and verifies that the page is loaded by checking if the link is visible.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyDiscordPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    if (await this.lnkDiscord.isHidden()) {
      await this.clickThenWait(this.btnMore)
      await this.clickThenWait(this.lnkDiscord)
    }
  }
}

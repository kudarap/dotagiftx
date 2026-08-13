import { expect, Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class HeroesPage extends CommonPageComponent {
  private readonly lnkHeroesPage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkHeroesPage = page.locator('a[href="/heroes"]')
  }

  /**
   * Clicks the Heroes link and verifies that the page is loaded by checking the page title.
   * @returns {Promise<void>} A promise that resolves when the verification is complete.
   */
  async clickAndVerifyHeroesPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    await this.clickThenWait(this.lnkHeroesPage)
    await expect(this.page).toHaveTitle(PageTitles.HEROES_PAGE)
  }
}

import { expect, Locator, Page } from '@playwright/test'
import { CommonPageComponent } from './common_page_component'
import { PageTitles } from '../constants/page_titles'

export class GuidesPage extends CommonPageComponent {
  private readonly lnkGuidesPage: Locator

  constructor(page: Page) {
    super(page)
    this.lnkGuidesPage = page.locator('a[href="/guides"]')
  }

  async clickAndVerifyGuidesPage(): Promise<void> {
    await this.waitForPageToLoad(this.page)
    if (await this.lnkGuidesPage.isHidden()) {
      await this.clickThenWait(this.btnMore)
      await this.clickThenWait(this.lnkGuidesPage)
    }
    await expect(this.page).toHaveTitle(PageTitles.GUIDES_PAGE)
  }
}

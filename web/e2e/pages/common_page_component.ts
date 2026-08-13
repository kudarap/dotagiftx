import { Page, Locator, expect, test } from '@playwright/test'
import { PlaywrightFunctions } from '../utilities/playwright_functions'
import { PageTitles } from '../constants/page_titles'

export class CommonPageComponent extends PlaywrightFunctions {
  private readonly gbWelcomeModal: Locator
  private readonly btnGotIt: Locator
  private readonly chkDontShowAgain: Locator
  public readonly btnMore: Locator

  constructor(page: Page) {
    super(page)
    this.gbWelcomeModal = page.locator("[aria-labelledby='welcome-dialog-title']")
    this.btnGotIt = page.getByRole('button', { name: 'Got it' })
    this.chkDontShowAgain = page.locator('.PrivateSwitchBase-input')
    this.btnMore = page.locator("//span[text()='More']")
  }

  /**
   * Closes the welcome modal if it is present on the page.
   * Waits for the page to load before checking for the modal's visibility.
   * If the modal is visible, clicks the close button and verifies that the modal is hidden.
   * @returns {Promise<void>} A promise that resolves when the modal is closed or if it was not present.
   */
  async closeWelcomeModalIfPresent(): Promise<void> {
    if (await this.isLocatorVisible(this.gbWelcomeModal)) {
      await this.chkDontShowAgain.check()
      await this.btnGotIt.click()
      await expect(this.gbWelcomeModal).toBeHidden()
    }
  }
}

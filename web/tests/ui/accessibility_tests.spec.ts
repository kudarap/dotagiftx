import { expect, test } from '@playwright/test'
import { DOTAGIFTX_URL } from '../../automation-files/constants/urls'
import { PageObjectModel } from '../../automation-files/pages/page_object_model'
import { CommonPageComponent } from '../../automation-files/pages/common_page_component'

test.describe('Accessibility Tests', () => {
  let commonPageComponent: CommonPageComponent
  let pom: PageObjectModel

  test.beforeEach(async ({ page }) => {
    commonPageComponent = new CommonPageComponent(page)
    pom = new PageObjectModel(page)
    await commonPageComponent.navigateToUrl(DOTAGIFTX_URL)
    await commonPageComponent.closeWelcomeModalIfPresent()
  })

  test('Verify Dota GiftX pages are loaded and accessible', async () => {
    await test.step('Verify Dota GiftX Home Page is loaded and accessible', async () => {
      await pom.dotaGiftXMainPage.clickAndVerifyDotaGiftXMainPageLoad()
    })

    await test.step('Verify Treasures Page is loaded and accessible', async () => {
      await pom.treasuresPage.clickAndVerifyTreasuresPage()
    })

    await test.step('Verify Heroes Page is loaded and accessible', async () => {
      await pom.heroesPage.clickAndVerifyHeroesPage()
    })

    await test.step('Verify Dota GiftX Plus Page is loaded and accessible', async () => {
      await pom.dotaGiftXPlus.clickAndVerifyDotaGiftxPlusPage()
    })

    await test.step('Verify Mobile Page is loaded and accessible', async () => {
      await pom.mobilePage.clickAndVerifyMobilePage()
    })

    await test.step('Verify Bans Page is loaded and accessible', async () => {
      await pom.bansPage.clickAndVerifyBansPage()
    })

    await test.step('Verify Guides Page is loaded and accessible', async () => {
      await pom.guidesPage.clickAndVerifyGuidesPage()
    })

    await test.step('Verify FAQs Page is loaded and accessible', async () => {
      await pom.faqsPage.clickAndVerifyFaqsPage()
    })

    await test.step('Verify Middle Man Page is loaded and accessible', async () => {
      await pom.middleManPage.clickAndVerifyMiddleManPage()
    })

    await test.step('Verify Moderators Page is loaded and accessible', async () => {
      await pom.moderatorsPage.clickAndVerifyModeratorsPage()
    })

    await test.step('Verify Updates Page is loaded and accessible', async () => {
      await pom.updatesPage.clickAndVerifyUpdatesPage()
    })
  })
})

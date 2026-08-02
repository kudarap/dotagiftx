import { test, expect } from '@playwright/test'
import { probeBackend } from './helpers'

let backendUp = true

test.beforeAll(async () => {
  backendUp = await probeBackend()
})

test.describe('home page', () => {
  test('renders header, search bar, trending and market stats', async ({ page }) => {
    test.skip(!backendUp, 'backend API is not reachable')

    const response = await page.goto('/')

    expect(response, 'GET /').not.toBeNull()
    expect(response!.ok(), 'GET / should return ok').toBeTruthy()

    await expect(page).toHaveTitle(/DotagiftX :: Dota 2 Giftables Community Market/)
    await expect(page.getByRole('link', { name: 'Treasures' })).toBeVisible()

    await expect(page.getByPlaceholder('Search for item name, hero, treasure')).toBeVisible()

    await expect(page.getByText('Trending')).toBeVisible()
    await expect(page.getByText('New Buy Orders')).toBeVisible()
    await expect(page.getByText('New Sell Listings')).toBeVisible()

    for (const stat of ['Available Offers', 'Buy Orders', 'On Reserved', 'Delivered Items']) {
      await expect(page.getByText(stat, { exact: true })).toBeVisible()
    }
  })

  test('submitting a search query navigates to the search page', async ({ page }) => {
    test.skip(!backendUp, 'backend API is not reachable')

    await page.goto('/')

    const input = page.getByPlaceholder('Search for item name, hero, treasure')
    await input.fill('Aghanim')
    await input.press('Enter')

    await page.waitForURL('**/search?q=Aghanim')

    await expect(page).toHaveTitle(/DotagiftX :: Search/)
  })
})

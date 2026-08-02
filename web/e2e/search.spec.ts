import { test, expect } from '@playwright/test'
import { probeBackend } from './helpers'

let backendUp = true

test.beforeAll(async () => {
  backendUp = await probeBackend()
})

test.describe('search page', () => {
  test('renders default search results', async ({ page }) => {
    test.skip(!backendUp, 'backend API is not reachable')

    const response = await page.goto('/search')

    expect(response, 'GET /search').not.toBeNull()
    expect(response!.ok(), 'GET /search should return ok').toBeTruthy()

    await expect(page).toHaveTitle(/DotagiftX :: Search/)
    await expect(page.locator('table tbody tr').first()).toBeVisible()
  })

  test('searching for a term shows the result count and matching items', async ({ page }) => {
    test.skip(!backendUp, 'backend API is not reachable')

    await page.goto('/search?q=Aghanim')

    await expect(page.getByText(/results for "Aghanim"/)).toBeVisible()
    await expect(page.locator('table tbody tr').first()).toBeVisible()
  })
})

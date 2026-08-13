import { test, expect } from '@playwright/test'
import { DOTAGIFTX_URL } from './constants/urls'

test('index', async ({ page }) => {
  await page.goto(`${DOTAGIFTX_URL}`)

  await expect(page).toHaveTitle(/Community Market/)
})

test('item page', async ({ page }) => {
  await page.goto(`${DOTAGIFTX_URL}/the-rat-king-chen`)

  await expect(page).toHaveTitle(/The Rat King/)
})

test('profile page', async ({ page }) => {
  await page.goto(`${DOTAGIFTX_URL}/profiles/76561198088587178`)

  await expect(page).toHaveTitle(/kudarap/)
})

import { test, expect } from '@playwright/test'

test.describe('navigation', () => {
  const headerLinks: Array<{ label: string; path: string }> = [
    { label: 'Bans', path: '/bans' },
    { label: 'Rules', path: '/rules' },
    { label: 'Mobile', path: '/download' },
    { label: 'Sign in', path: '/login' },
  ]

  for (const { label, path } of headerLinks) {
    test(`header link "${label}" navigates to ${path}`, async ({ page }) => {
      await page.goto('/')

      await page.getByRole('link', { name: label, exact: true }).click()
      await page.waitForURL(`**${path}`)

      expect(page.url()).toContain(path)
      await expect(page).toHaveTitle(/DotagiftX/)
    })
  }

  test('header "More" menu exposes sub-pages', async ({ page }) => {
    await page.goto('/')

    await page.getByText('More', { exact: true }).hover()
    await page.getByRole('menuitem', { name: 'FAQs' }).click()
    await page.waitForURL('**/faqs')

    await expect(page.getByRole('heading', { name: 'Frequently Asked Questions' })).toBeVisible()
  })

  test('footer links navigate to their pages', async ({ page }) => {
    await page.goto('/')

    await page.getByRole('link', { name: 'Privacy' }).click()
    await page.waitForURL('**/privacy')
    await expect(page.getByRole('heading', { name: 'Privacy Policy' })).toBeVisible()

    await page.getByRole('link', { name: 'Donate' }).click()
    await page.waitForURL('**/donate')
    await expect(page.getByRole('heading', { name: 'How can I donate?' })).toBeVisible()
  })
})

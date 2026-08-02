import { test, expect } from '@playwright/test'

test.describe('login page', () => {
  test('shows Steam sign in options and security notice', async ({ page }) => {
    const response = await page.goto('/login')

    expect(response, 'GET /login').not.toBeNull()
    expect(response!.ok(), 'GET /login should return ok').toBeTruthy()

    await expect(page).toHaveTitle(/DotagiftX :: Sign In/)
    await expect(page.getByRole('heading', { name: /Signing in to/ })).toBeVisible()
    await expect(page.getByRole('link', { name: /Sign in through Steam/ })).toBeVisible()
    await expect(
      page.getByText('This website is not affiliated with Valve Corporation or Steam.')
    ).toBeVisible()
  })
})

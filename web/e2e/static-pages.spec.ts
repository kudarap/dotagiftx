import { test, expect } from '@playwright/test'

const IGNORED_CONSOLE_ERRORS = [
  'cannot contain a nested',
  'cannot be a descendant of',
  'This will cause a hydration error',
  'validateDOMNesting',
  'Download the React DevTools',
  'hydration error',
]

function collectPageErrors(page: import('@playwright/test').Page): () => string[] {
  const errors: string[] = []
  page.on('pageerror', err => errors.push(err.message))
  page.on('console', msg => {
    if (msg.type() !== 'error') return
    if (IGNORED_CONSOLE_ERRORS.some(prefix => msg.text().includes(prefix))) return
    errors.push(msg.text())
  })
  return () => errors
}

const staticPages: Array<{
  path: string
  heading: string
  title: string
  dataDependent?: boolean
}> = [
  { path: '/faqs', heading: 'Frequently Asked Questions', title: 'DotagiftX' },
  { path: '/rules', heading: 'Rules', title: 'DotagiftX' },
  { path: '/privacy', heading: 'Privacy Policy', title: 'DotagiftX' },
  { path: '/bans', heading: 'Banned users', title: 'DotagiftX' },
  { path: '/donate', heading: 'How can I donate?', title: 'DotagiftX' },
  { path: '/guides', heading: 'Guides & Tips', title: 'DotagiftX' },
  { path: '/middleman', heading: 'Middleman', title: 'DotagiftX' },
  { path: '/moderators', heading: 'Moderators', title: 'DotagiftX', dataDependent: true },
  { path: '/updates', heading: 'Updates', title: 'DotagiftX' },
  { path: '/download', heading: 'DotagiftX for Mobile', title: 'DotagiftX' },
  { path: '/plus', heading: 'Dotagift Plus', title: 'DotagiftX' },
]

for (const { path, heading, title, dataDependent } of staticPages) {
  test.describe(`static page: ${path}`, () => {
    test('renders with expected title, heading, header and footer', async ({ page }) => {
      const getErrors = collectPageErrors(page)

      const response = await page.goto(path)

      expect(response, `GET ${path}`).not.toBeNull()
      expect(response!.ok(), `GET ${path} should return ok`).toBeTruthy()

      await expect(page).toHaveTitle(new RegExp(title))
      await expect(page.getByRole('heading', { level: 1, name: new RegExp(heading) })).toBeVisible()

      // Shared header and footer chrome.
      await expect(page.getByRole('link', { name: 'Treasures' })).toBeVisible()
      await expect(page.getByRole('link', { name: 'Sign in', exact: true })).toBeVisible()
      await expect(page.getByText('not affiliated with Valve or Steam')).toBeVisible()

      // Data-dependent pages may legitimately log fetch errors when the
      // referenced records are absent from the environment.
      if (!dataDependent) {
        expect(getErrors(), 'no page errors').toEqual([])
      }
    })
  })
}

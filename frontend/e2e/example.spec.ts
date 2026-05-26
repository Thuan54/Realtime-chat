import { test, expect } from '@playwright/test'

test.describe('Chat Application Smoke Journey', () => {
  test('user can login, join channel, send/receive message', async ({ page }) => {
    // 1. Navigate to app
    await page.goto('/')
    
    // 2. Login with seeded test user (credentials from 07-testing-qa.md fixtures)
    await page.getByLabel('Email').fill('testuser@example.com')
    await page.getByLabel('Password').fill('TestPass123!')
    await page.getByRole('button', { name: 'Log In' }).click()
    
    // 3. Verify channel list loads (FR-03)
    await expect(page.getByRole('navigation', { name: 'Channels' })).toBeVisible()
    
    // 4. Join a public channel
    await page.getByRole('button', { name: 'Join Channel' }).click()
    await page.getByPlaceholder('Channel name').fill('general')
    await page.getByRole('button', { name: 'Join' }).click()
    
    // 5. Send a message (FR-05)
    const messageText = `E2E test message ${Date.now()}`
    await page.getByPlaceholder('Type a message...').fill(messageText)
    await page.keyboard.press('Enter')
    
    // 6. Verify message appears in history (FR-06)
    await expect(page.getByText(messageText)).toBeVisible()
    
    // 7. Verify real-time receipt via WebSocket (NFR-01)
    // Note: Full WS assertion requires mocking; covered in integration layer
  })

  test.fixme('offline queue sync after reconnect (FR-08)', async ({ page }) => {
    // Implementation requires Playwright network interception + WS mock
    // Deferred to integration test layer per 07-testing-qa.md
    await page.goto('/')
    // ... future implementation here
  })
})
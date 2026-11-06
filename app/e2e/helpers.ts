import { expect, type APIRequestContext, type Page } from "@playwright/test";

export const PASSWORD = process.env.SEED_PASSWORD || "password123";
export const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

export async function loginViaUI(page: Page, identifier: string) {
  await page.goto("/login");
  await page.locator("#identifier").fill(identifier);
  await page.locator("#password").fill(PASSWORD);
  await page.locator('button[type="submit"]').click();
  // first visit compiles the dashboard route in dev mode, so allow extra time
  await expect(page).toHaveURL(/\/dashboard$/, { timeout: 30_000 });
}

export async function loginViaAPI(request: APIRequestContext, identifier: string) {
  const res = await request.post(`${BACKEND_URL}/auth/login`, {
    data: { identifier, password: PASSWORD, turnstileToken: "e2e" },
  });
  expect(res.status(), `login as ${identifier}`).toBe(202);
}

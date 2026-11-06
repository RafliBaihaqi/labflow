import { expect, test } from "@playwright/test";
import { BACKEND_URL, loginViaAPI, loginViaUI } from "./helpers";

test("customer signs in through the UI", async ({ page }) => {
  await loginViaUI(page, "customer");
  await page.reload();
  await expect(page).toHaveURL(/\/dashboard$/); // session survives a reload
});

test("wrong password is rejected", async ({ page }) => {
  await page.goto("/login");
  await page.locator("#identifier").fill("customer");
  await page.locator("#password").fill("not-the-password");
  await page.locator('button[type="submit"]').click();
  await expect(page).toHaveURL(/\/login$/);
});

test.describe("RBAC on the API", () => {
  test("customer cannot list quotations", async ({ request }) => {
    await loginViaAPI(request, "customer");
    const res = await request.get(`${BACKEND_URL}/api/quotations`);
    expect(res.status()).toBe(403);
  });

  test("admin can list quotations", async ({ request }) => {
    await loginViaAPI(request, "admin");
    const res = await request.get(`${BACKEND_URL}/api/quotations`);
    expect(res.status()).toBe(200);
  });

  test("anonymous requests are unauthorized", async ({ request }) => {
    const res = await request.get(`${BACKEND_URL}/api/reservations`);
    expect(res.status()).toBe(401);
  });

  test("logout revokes the session token", async ({ playwright }) => {
    const session = await playwright.request.newContext();
    await loginViaAPI(session, "micro_qc");
    const { cookies } = await session.storageState();
    const token = cookies.find((c) => c.name === "Authorization")?.value;
    expect(token).toBeTruthy();
    const withOldToken = { headers: { Cookie: `Authorization=${token}` } };
    const fresh = await playwright.request.newContext();
    expect((await fresh.get(`${BACKEND_URL}/api/auth/me`, withOldToken)).status()).toBe(200);
    await fresh.dispose();

    expect((await session.post(`${BACKEND_URL}/auth/logout`)).ok()).toBeTruthy();

    // Re-send the old token explicitly: logout clears the cookie, but the token
    // itself must also be rejected server-side.
    const replay = await playwright.request.newContext();
    expect((await replay.get(`${BACKEND_URL}/api/auth/me`, withOldToken)).status()).toBe(401);
    await Promise.all([session.dispose(), replay.dispose()]);
  });
});

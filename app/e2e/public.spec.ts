import { expect, test } from "@playwright/test";

test("landing page renders and accepts a quotation request", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { level: 1 })).toContainText("sample to certificate");

  await page.getByPlaceholder("Email").fill("buyer@example.com");
  await page.getByPlaceholder("Company name").fill("Example Labs");
  await page.getByPlaceholder("What would you like tested?").fill("Bioburden on 3 batches");
  await page.getByRole("button", { name: "Send request" }).click();

  await expect(page.getByText("Request sent")).toBeVisible();
});

test("protected pages redirect anonymous visitors to login", async ({ page }) => {
  await page.goto("/dashboard");
  await expect(page).toHaveURL(/\/login$/);
});
